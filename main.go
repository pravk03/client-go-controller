package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	networking "k8s.io/api/networking/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var kubeconfig string

func SetupSignalHandler() context.Context {

	ctx, cancel := context.WithCancel(context.Background())

	// Handle signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan // Wait for a signal
		fmt.Println("Received termination signal, shutting down...")
		cancel() // Cancel the context to signal shutdown
	}()

	return ctx

}

func main() {

	klog.InitFlags(nil)
	flag.Parse()

	ctx := SetupSignalHandler()
	logger := klog.FromContext(ctx)
	// Create Kubernetes clientset

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logger.Info("Error building kubeconfig")
		klog.FlushAndExit(klog.ExitFlushTimeout, 1)
	}

	clientset, _ := kubernetes.NewForConfig(cfg)

	// Create Informer for network policies
	kubeInformerFactory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	networkPolicyInformer := kubeInformerFactory.Networking().V1().NetworkPolicies().Informer()

	logger.Info("Starting controller")

	// Event Handler: When a Network policy is added, deleted or updated
	networkPolicyInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			// TODO: Handle network policy update
			metadata, _ := newObj.(*networking.NetworkPolicy)
			logger.Info("Network Policy updated", "name", metadata.Name, "namespace", metadata.Namespace)
		},
		AddFunc: func(obj interface{}) {
			// TODO: Handle network policy add
			metadata, _ := obj.(*networking.NetworkPolicy)
			logger.Info("Network Policy added", "name", metadata.Name, "namespace", metadata.Namespace)
		},
		DeleteFunc: func(obj interface{}) {
			// TODO: Handle network policy delete
			metadata, _ := obj.(*networking.NetworkPolicy)
			logger.Info("Network Policy deleted", "name", metadata.Name, "namespace", metadata.Namespace)
		},
	})

	// Start the informer and run forever
	kubeInformerFactory.Start(ctx.Done())
	kubeInformerFactory.WaitForCacheSync(ctx.Done())

	<-ctx.Done()
}

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "Path to kubeconfig. Only required if running outside the cluster")
}

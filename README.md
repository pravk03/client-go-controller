# client-go-controller


## Create a kind cluster

To create a Kubernetes cluster using Kind (Kubernetes in Docker) with 2 nodes, you can use the following command:

```
    kind create cluster --name my-cluster --config sample_artifacts/kind_cluster.yaml
```

This command creates a Kind cluster named "my-cluster" with one control-plane node and one worker node. Adjust the --name flag and other parameters as needed for your specific requirements.

## Runnnig the controller from outside the cluster

### Build 

```
    make build
```

### Run the controller

```
    bin/main -kubeconfig=$HOME/.kube/config
```

### Verify

Add/delete/update network policy in a separate terminal and check log messages on the terminal running the controler

```
    kubectl apply -f sample_artifacts/dummy_network_policy.yaml
    // LOG: "Network Policy added" name="allow-all-traffic" namespace="default"

    kubectl delete networkpolicies allow-all-traffic
    // LOG: "Network Policy deleted" name="allow-all-traffic" namespace="default"

```


## Run the controller with the cluster 

TODO : create a docker container and add a pod to run the controller with a pod.
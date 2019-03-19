# Nemo for Kubernetes
Nemo come bundled with [Dory](https://github.com/hpe-storage/dory) for Kubernetes. Dory is a FlexVolume driver for Docker Volume plugins. Dory also comes with an out-of-tree dynamic provisioner, Doryd. Dory and Doryd is the Open Source projects incorporated into the HPE Nimble Storage container integration and works just like the "The Real Thing".

Nemo will run as a DaemonSet and Doryd as a Deployment. Dory the FlexVolume driver is part of the DaemonSet. The DaemonSet itself **requires Kubernetes 1.10+** due to its reliance on mount propagation.

## Installing
The following specs should be deployed as a cluster admin:
```
kubectl create -f https://raw.githubusercontent.com/NimbleStorage/Nemo/master/runtime/k8s/daemonset-nemod.yaml
kubectl create -f https://raw.githubusercontent.com/NimbleStorage/Nemo/master/runtime/k8s/deploy-doryd.yaml
```
It's now possible to create StorageClasses and Persistent Volume Claims or use the FlexVolume driver inline. The provisioner listens on `dev.hpe.com/nemo`.

### GKE workaround ###
When deploying `doryd` on GKE, a custom role binding for your particular user running `kubectl` needs to be created:
```
kubectl create clusterrolebinding cluster-admin-nemo --clusterrole=cluster-admin --user=user@fqdn.com
```
Also, the FlexVolume path on GKE is in a non-standard location. A sample DaemonSet has been provided as such:
```
kubectl create -f https://raw.githubusercontent.com/NimbleStorage/Nemo/master/runtime/k8s/daemonset-nemod-gke.yaml
kubectl create -f https://raw.githubusercontent.com/NimbleStorage/Nemo/master/runtime/k8s/deploy-doryd-gke.yaml
```

## Using
The nature of Nemo, being a node local affair, it's only advisable to only deploy Pods or StatefulSets as they are sticky to the node that they are deployed on.

An example StorageClass and StatefulSet:
```
kubectl create -f https://raw.githubusercontent.com/NimbleStorage/Nemo/master/runtime/k8s/sc-transactionaldb.yaml
kubectl create -f https://raw.githubusercontent.com/NimbleStorage/Nemo/master/runtime/k8s/statefulset-mariadb.yaml
```

# Customization
Nemo for Kubernetes uses the same parameters as the standalone plugin. The DaemonSet will create a new pool on a loopback device if no pool name is passed through the environment variables. 

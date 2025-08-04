# Virtual-Machines-Management
A Kubernetes native application which can manage Virtual Machines'S manipulation.

## Features
- Create Virtual Machines
- Delete Virtual Machines
- List Virtual Machines
- Get Virtual Machines details

## How to run
### Custom API Server
#### Change directory to custom-apiserver
```bash
cd custom-apiserver
```
#### Build the Binary
```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o artifacts/simple-image/virtual-machines-apiserver
```

#### Build the Docker Image
```bash
docker build -t MYPREFIX/virtual-machines-apiserver:MYTAG ./artifacts/simple-image
docker push MYPREFIX/virtual-machines-apiserver:MYTAG
```
Where `MYPREFIX` is your Docker Hub username and `MYTAG` is the tag you want to use.

#### Deploy the Custom API Server
```bash
kubectl apply -f artifacts/example/ns.yaml
kubectl apply -f artifacts/example
```

#### Verify the Custom API Server is running
Test if the custom API server is running by checking the apiservices and the pods in the `kube-system` namespace:

```bash
kubectl get apiservices -n kube-system
kubectl get pods -n vms
```

### Custom Controller
#### Change directory to custom-controller
```bash
cd custom-controller
```

#### Build the image
```bash
make docker-build IMG=MYPREFIX/virtual-machines-controller:MYTAG
```
Where `MYPREFIX` is your Docker Hub username and `MYTAG` is the tag you want to use.

#### Push the image
```bash
docker push MYPREFIX/virtual-machines-controller:MYTAG
```

#### Deploy the Custom Controller
1. Create a Kubernetes secret manifest for your AWS credentials, then apply it:
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: aws-credentials
  namespace: custom-controller-system
type: Opaque
stringData:
  access_key: <YOUR_AWS_ACCESS_KEY>
  secret_key: <YOUR_AWS_SECRET_KEY>
```

```bash
kubectl apply -f aws-credentials.yaml
```

2. Locate your operator’s deployment YAML file, typically in config/manager/manager.yaml. Change the image field to point to your custom controller image:

```yaml
spec:
  containers:
    - name: manager
      image: MYPREFIX/virtual-machines-controller:MYTAG
```

Then apply the deployment:

```bash
make deploy IMG=MYPREFIX/virtual-machines-controller:MYTAG
```

3. Apply the RBAC rules:
```bash
kubectl apply -f virtualmachine-operator-permissions.yaml
```

#### Verify the Custom Controller is running
Check if the custom controller is running by checking the pods in the `custom-controller-system` namespace

```bash
kubectl get pods -n custom-controller-system
```

### Create a Virtual Machine
To create a Virtual Machine, apply the following manifest:
```bash
kubectl apply -f custom-apiserver/artifacts/virtualmachine/02-virtualmachine.yaml
```

### List Virtual Machines
To list all Virtual Machines, use the following command:
```bash
kubectl get virtualmachines
```

### Get Virtual Machine details
To get the details of a specific Virtual Machine, use:
```bash
kubectl get virtualmachine <vm-name> -o yaml
```

### Delete a Virtual Machine
To delete a Virtual Machine, use the following command:
```bash
kubectl delete virtualmachine <vm-name>
```
This will remove the Virtual Machine from the cluster.
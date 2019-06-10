### Prerequisites
In order to build and deploy the project using the 
instructions in this README, make sure you have 
the following prerequisites installed locally  
* minikube
* make >= 3.82

#### Build
You can build a docker image of the external metrics API server using the following command:
 
`make build`

#### Deploy
Once the docker image is built, you can deploy it using:

`make deploy`

#### Access the external metrics
Access the external metrics by querying the kubernetes API directly:

````kubectl get --raw /apis/external.metrics.k8s.io/v1beta1/namespaces/default/tube-default````

### Tear down
When you are finished, you can clean up by executing

`make delete`
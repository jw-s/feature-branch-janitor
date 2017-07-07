package k8sclient

import (
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	ClientSet *kubernetes.Clientset
}

func New() Client {

	config, err := rest.InClusterConfig()

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {

		panic(err.Error())
	}
	return Client{
		ClientSet: clientset,
	}

}
func (client *Client) GetDeploymentsWithAnnotation(namespace, annotationName, annotationValue string) []v1beta1.Deployment {

	deploymentList, err := client.ClientSet.ExtensionsV1beta1().Deployments(namespace).List(v1.ListOptions{})

	if err != nil {

		panic(err.Error())
	}

	var deployments []v1beta1.Deployment

	for _, deployment := range deploymentList.Items {

		annotations := deployment.ObjectMeta.Annotations

		if v, ok := annotations[annotationName]; ok {
			if v == annotationValue {
				deployments = append(deployments, deployment)
			}
		}

	}

	return deployments
}

package k8s

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Client struct {
	ClientSet *kubernetes.Clientset
}

type SecretData map[string][]byte

type DeploymentsInNamespaces map[string][]*v1beta1.Deployment

func New() (*Client, error) {

	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)

	if err != nil {

		return nil, err
	}

	return &Client{
		ClientSet: clientset,
	}, err

}
func (client *Client) GetDeploymentsWithAnnotations(namespaces []string, annotationName ...string) DeploymentsInNamespaces {

	namespaceNames := GetNamespaceNames(client.GetNamespaces(namespaces))

	deploymentsInNamespaces := make(DeploymentsInNamespaces)

	for _, namespace := range namespaceNames {

		var deployments []*v1beta1.Deployment

		deploymentList, err := client.ClientSet.ExtensionsV1beta1().Deployments(namespace).List(meta.ListOptions{})

		if err != nil {

			log.Fatal(err.Error())
		}

		for _, deployment := range deploymentList.Items {

			foundDeployment := false

			annotations := deployment.ObjectMeta.Annotations

			for _, annotation := range annotationName {
				if _, ok := annotations[annotation]; ok {
					foundDeployment = true
				} else {
					foundDeployment = false
				}

			}

			if foundDeployment {
				deployments = append(deployments, &deployment)
			}

		}
		deploymentsInNamespaces[namespace] = deployments
	}

	return deploymentsInNamespaces
}

func (client *Client) DeleteDeployments(deploymentsInNamespaces DeploymentsInNamespaces) error {

	for namespace, deployments := range deploymentsInNamespaces {

		if len(deployments) == 0 {
			continue
		}

		if err := client.deleteDeployments(namespace, deployments...); err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) deleteDeployments(namespace string, deployments ...*v1beta1.Deployment) error {

	deploymentsClient := client.ClientSet.AppsV1beta1().Deployments(namespace)

	deletePolicy := meta.DeletePropagationForeground

	deploymentNames := GetDeploymentNames(deployments)

	for _, name := range deploymentNames {

		if err := deploymentsClient.Delete(name, &meta.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		}); err != nil {
			return err
		}

	}

	log.Infof("Namespace[%s]: deleted the following deployments %s", namespace, deploymentNames)

	return nil
}

func (client *Client) GetNamespaces(namespaces []string) []v1.Namespace {

	var namespaceList []v1.Namespace

	if len(namespaces) == 0 {
		list, err := client.ClientSet.Namespaces().List(meta.ListOptions{})

		if err != nil {
			return namespaceList
		}

		return list.Items

	}

	for _, namespace := range namespaces {
		n, err := client.ClientSet.Namespaces().Get(namespace, meta.GetOptions{})

		if err != nil {
			continue
		}
		namespaceList = append(namespaceList, *n)
	}

	return namespaceList
}

func GetNamespaceNames(namespaces []v1.Namespace) []string {
	var namespaceNames []string

	for _, namespace := range namespaces {
		namespaceNames = append(namespaceNames, namespace.GetName())
	}

	return namespaceNames
}

func GetDeploymentNames(deployments []*v1beta1.Deployment) []string {
	var deploymentNames []string

	for _, deployment := range deployments {
		deploymentNames = append(deploymentNames, deployment.GetName())
	}

	return deploymentNames
}
func (client *Client) GetSecret(namespace, secretName string) SecretData {

	secret, err := client.ClientSet.CoreV1().Secrets(namespace).Get(secretName, meta.GetOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	return secret.Data

}

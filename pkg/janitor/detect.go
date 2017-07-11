package janitor

import (
	"context"

	"github.com/JoelW-S/feature-branch-janitor/pkg/k8s"
)

func (janitor *Janitor) DeleteDeploymentsWithDeletedBranches() error {

	deployments := janitor.GetDeploymentsWithDeletedBranches()

	if err := janitor.K8sClient.DeleteDeployments(deployments); err != nil {
		return err
	}

	return nil

}

func (janitor *Janitor) GetDeploymentsWithDeletedBranches() k8s.DeploymentsInNamespaces {

	deploymentsInNamespaces := janitor.K8sClient.GetDeploymentsWithAnnotations(janitor.Namespaces, janitor.BranchAnnotation, janitor.RepositoryAnnotation)

	for namespace, deployments := range deploymentsInNamespaces {

		for i, deployment := range deployments {

			repo := deployment.ObjectMeta.Annotations[janitor.RepositoryAnnotation]

			branchName := deployment.ObjectMeta.Annotations[janitor.BranchAnnotation]

			if _, found := janitor.GithubClient.GetBranch(context.Background(), janitor.Owner, repo, branchName); found {
				deploymentsInNamespaces[namespace] = append(deployments[:i], deployments[i+1:]...)
			}

		}

	}

	return deploymentsInNamespaces

}

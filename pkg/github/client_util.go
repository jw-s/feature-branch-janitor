package poller

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ListBranches(owner, repository string) []string {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: secretToToken()},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	branches, err := client.Repositories.listBranches(ctx, owner, repository, &ListOptions{})

	if err != nil {

		panic(err.Error())

	}

	return getNamesFromBranches(branches)

}

func getNamesFromBranches(branches []*Branch) {

	branchNames := make([]string, len(branches))

	for branch := range branches {

		branchNames := append(branchNames, branch.GetName())

	}
	return branchNames
}

func secretToToken() string {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	secret, err := clientset.CoreV1().Secrets("cleanup").Get("bitbucket-dof", meta_v1.GetOptions{}).String()
	if err != nil {
		panic(err.Error())
	}

	return secret

}

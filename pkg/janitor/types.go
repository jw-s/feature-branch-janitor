package janitor

import (
	"os"
	"time"

	"github.com/JoelW-S/feature-branch-janitor/pkg/github"

	"github.com/JoelW-S/feature-branch-janitor/pkg/k8s"

	log "github.com/sirupsen/logrus"
)

type ticker interface {
	Tick() <-chan time.Time
	Stop()
}

type Cycle struct {
	*time.Ticker
}

func (c *Cycle) Tick() <-chan time.Time { return c.C }

func (c *Cycle) Stop() { c.Ticker.Stop() }

func NewCycle(duration time.Duration) *Cycle {
	return &Cycle{time.NewTicker(duration)}
}

type Janitor struct {
	Cycle                *Cycle
	Namespaces           []string
	BranchAnnotation     string
	RepositoryAnnotation string
	Owner                string
	GithubClient         *github.Client
	K8sClient            *k8s.Client
}

func NewJanitor(cycle *Cycle,
	namespaces []string,
	branchAnnotation string,
	repositoryAnnotation string,
	owner string,
	githubClient *github.Client,
	k8sClient *k8s.Client) *Janitor {

	return &Janitor{
		Cycle:                cycle,
		Namespaces:           namespaces,
		BranchAnnotation:     branchAnnotation,
		RepositoryAnnotation: repositoryAnnotation,
		Owner:                owner,
		GithubClient:         githubClient,
		K8sClient:            k8sClient,
	}
}

func (janitor *Janitor) Roam(exitChan <-chan (os.Signal)) {

	cycle := janitor.Cycle

	cycles := 1

	for {
		select {
		case <-cycle.Tick():

			log.WithFields(log.Fields{
				"cycle": cycles,
			}).Info("Starting janitor duties")

			if err := janitor.DeleteDeploymentsWithDeletedBranches(); err != nil {
				log.WithFields(log.Fields{
					"cleanup": false,
				}).Fatal("Cleanup failed")
			}

			log.WithFields(log.Fields{
				"cycle":   cycles,
				"cleanup": true,
			}).Info("Cleanup cycle complete")

			cycles++

		case <-exitChan:

			log.Infof("received signal: %v", <-exitChan)
			cycle.Stop()
			os.Exit(1)

		}
	}
}

package main

import (
	"os"
	"os/signal"

	"github.com/urfave/cli"

	"time"

	"github.com/JoelW-S/feature-branch-janitor/pkg/github"
	"github.com/JoelW-S/feature-branch-janitor/pkg/janitor"
	"github.com/JoelW-S/feature-branch-janitor/pkg/k8s"
	log "github.com/sirupsen/logrus"
)

var (
	appName         = "feature branch janitor"
	usage           = "A Kubernetes operator used to cleanup feature branch deployments when the branch has been merged and deleted."
	version         string
	owner           string
	branch          string
	repository      string
	secret          string
	secretNamespace string
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.JSONFormatter{})

}
func main() {

	app := cli.NewApp()

	app.Name = appName
	app.Version = version
	app.Usage = usage
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "owner, o",
			Value:       "JoelW-S",
			Usage:       "Scm owner, `Source control user account`",
			Destination: &owner,
		},
		cli.StringFlag{
			Name:        "branch-annotation, ba",
			Value:       "autodelete_branch",
			Usage:       "Annotation containing branch to watch",
			Destination: &branch,
		},
		cli.StringFlag{
			Name:        "repository-annotation, ra",
			Value:       "autodelete_repo",
			Usage:       "Annotation containg repo to watch",
			Destination: &repository,
		},
		cli.StringFlag{
			Name:        "secret, s",
			Value:       "repo-credentials",
			Usage:       "Kubernetes secret used to login to specified scm",
			Destination: &secret,
		},
		cli.StringFlag{
			Name:        "secret-namespace, sn",
			Value:       "cleanup",
			Usage:       "Namespace containing secret to login to specified scm",
			Destination: &secretNamespace,
		},
		cli.BoolFlag{
			Name:  "all-namespaces",
			Usage: "Watch all namespaces",
		},
		cli.StringSliceFlag{
			Name:  "namespaces",
			Usage: "Namespaces to watch",
		},
		cli.DurationFlag{
			Name: "duration, d", Value: time.Minute,
		},
	}

	app.Action = func(c *cli.Context) error {

		var namespaces []string

		if len(c.StringSlice("namespaces")) == 0 && !c.Bool("all-namespaces") {
			log.Fatal("Must have atleast one namespace set for watching")
		} else if !c.Bool("all-namespaces") {
			namespaces = c.StringSlice("namespaces")
		} else {
			namespaces = make([]string, 0)
		}

		client, err := k8s.New()

		if err != nil {
			log.WithFields(log.Fields{
				"not_in_kubernetes": true,
			}).Fatal("Application is not deployed in Kubernetes")
		}

		secretData := client.GetSecret(secretNamespace, secret)

		token, ok := secretData["token"]

		if !ok {
			log.WithFields(log.Fields{
				"exists": false,
			}).Fatal("Token retrieval issues")
		}

		myJanitor := janitor.NewJanitor(
			janitor.NewCycle(c.Duration("duration")),
			namespaces,
			branch,
			repository,
			owner,
			github.NewAuthenticatedClient(string(token)),
			client)

		sigChan := make(chan os.Signal, 1)

		defer close(sigChan)

		signal.Notify(sigChan)

		myJanitor.Roam(sigChan)

		return nil
	}

	app.Run(os.Args)

}

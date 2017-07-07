package main

import (
	"os"
	"os/signal"

	"github.com/JoelW-S/feature-branch-janitor/pkg/k8sclient"
	"github.com/Sirupsen/logrus"
)

func main() {

	c := make(chan os.Signal, 1)
	signal.Notify(c)
	go func() {
		logrus.Infof("received signal: %v", <-c)
		os.Exit(1)
	}()

	client := k8sclient.New()

	deployments := client.GetDeploymentsWithAnnotation("test", "test", "test")

	for _, deployment := range deployments {

		logrus.Info(deployment.Name)

	}
}

package main

import (
	"go.temporal.io/server/common/log"

	"github.com/CMSGov/cms-cloud-account-onboarding/app"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewDevelopment()
	c, err := client.NewClient(client.Options{
		Logger: log.NewZapAdapter(logger),
	})
	if err != nil {
		logger.Info("There was a problem with the logger")
	}

	defer c.Close()

	w := worker.New(c, app.AccountCreateQueue, worker.Options{})
	w.RegisterWorkflow(app.AccountOnboardingWorkflow)
	w.RegisterActivity(app.CreateAccount)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		logger.Info("There was a problem running the activity")
	}
}

package app

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func AccountOnboardingWorkflow(ctx workflow.Context, accountDetails CreateAccountInput) (string, error) {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
	}
	logger := workflow.GetLogger(ctx)

	ctx = workflow.WithActivityOptions(ctx, options)
	var result string
	logger.Info("Starting Workflow Activity...")
	err := workflow.ExecuteActivity(ctx, CreateAccount, accountDetails).Get(ctx, &result)
	return result, err
}

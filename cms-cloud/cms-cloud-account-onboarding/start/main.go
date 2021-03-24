package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/CMSGov/cms-cloud-account-onboarding/app"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Println(err)
	}

	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "account-onboarding-workflow",
		TaskQueue: app.AccountCreateQueue,
	}

	var (
		idx          = os.Args[1]
		ouName       = os.Args[2]
		accountName  = fmt.Sprintf("aws-sso-poc-%s", idx)
		accountEmail = fmt.Sprintf("aws-cms-oit-iusg-acct%s@gdit.com", idx)
		ssoEmail     = "faraaz@samtek-inc.com"
		firstName    = "Faraaz"
		lastName     = "Khan"
	)

	createAccountInput := app.CreateAccountInput{
		AccountName:          accountName,
		AccountEmail:         accountEmail,
		AccountSSOEmail:      ssoEmail,
		AccountSSOFirstName:  firstName,
		AccountSSOLastName:   lastName,
		OrganizationalUnitID: ouName,
	}

	we, err := c.ExecuteWorkflow(context.Background(), options, app.AccountOnboardingWorkflow, createAccountInput)
	if err != nil {
		log.Println(err)
	}

	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Println(err)
	}
	printResults(result, we.GetID(), we.GetRunID())
}

func printResults(result string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", result)
}

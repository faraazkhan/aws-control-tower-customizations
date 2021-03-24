package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin"
	"github.com/aws/aws-sdk-go-v2/service/ssoadmin/types"
)

var (
	region                 = os.Getenv("SSO_REGION")
	instance               = os.Getenv("SSO_INSTANCE_ARN")
	readOnlyGroupsList     = os.Getenv("GSS_READONLY_GROUPS")
	powerUserGroupsList    = os.Getenv("GSS_POWERUSER_GROUPS")
	adminGroupsList        = os.Getenv("GSS_ADMIN_GROUPS")
	adminPermissionSet     = os.Getenv("GSS_POWERUSER_PERMISSION_ARN")
	powerUserPermissionSet = os.Getenv("GSS_ADMIN_PERMISSION_ARN")
	readOnlyPermissionSet  = os.Getenv("GSS_READONLY_PERMISSION_ARN")
)

func createGSSAccountAssignments(accountID string) {

	groupsAndPermissions := map[string]string{
		adminPermissionSet:     adminGroupsList,
		powerUserPermissionSet: powerUserGroupsList,
		readOnlyPermissionSet:  readOnlyGroupsList,
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	svc := ssoadmin.NewFromConfig(cfg)

	for permissionSet, groups := range groupsAndPermissions {
		for _, group := range strings.Split(groups, ",") {
			_, err := svc.CreateAccountAssignment(context.TODO(), &ssoadmin.CreateAccountAssignmentInput{
				InstanceArn:      &instance,
				PermissionSetArn: &permissionSet,
				PrincipalType:    types.PrincipalTypeGroup,
				PrincipalId:      &group,
				TargetType:       types.TargetTypeAwsAccount,
				TargetId:         &accountID,
			})
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Account Assignment Created for Group: %s with Permission Set: %s\n", group, permissionSet)
		}
	}

}

// HandleEvent  parses event JSON for account details and passes them to createGSSAccountAssignments
func HandleEvent(ctx context.Context, accountCreationEvent events.SQSEvent) {
	for _, message := range accountCreationEvent.Records {
		var (
			managedAccountStatus ManagedAccountStatus
			eventBridgeEvent     EventBridgeEvent
		)

		err := json.Unmarshal([]byte(message.Body), &eventBridgeEvent)
		if err != nil {
			log.Fatal("Could not unmarshal Unquoted Body: ", err)
		}

		eventDetails := eventBridgeEvent.Detail

		if eventDetails.ServiceEventDetails.CreateManagedAccountStatus.State != "" {
			managedAccountStatus = eventDetails.ServiceEventDetails.CreateManagedAccountStatus
		} else if eventDetails.ServiceEventDetails.UpdateManagedAccountStatus.State != "" {
			managedAccountStatus = eventDetails.ServiceEventDetails.UpdateManagedAccountStatus
		} else {
			log.Fatal("Could not parse ManagedAccountStatus from Event")
		}

		var (
			eventAccountID   = managedAccountStatus.Account.AccountID
			eventAccountName = managedAccountStatus.Account.AccountName
			eventOUID        = managedAccountStatus.OrganizationalUnit.OrganizationalUnitID
			eventOUName      = managedAccountStatus.OrganizationalUnit.OrganizationalUnitName
		)

		if managedAccountStatus.State != "SUCCEEDED" {
			log.Fatal("Control Tower could not enroll account: ", eventAccountID)
		}

		log.Printf("Account Assignment Creation Started for Account Name: %s(%s) in OU: %s(%s)",
			eventAccountName, eventAccountID, eventOUName, eventOUID)
		createGSSAccountAssignments(eventAccountID)

	}

}

func main() {
	lambda.Start(HandleEvent)
}

package app

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalog/types"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"go.temporal.io/sdk/activity"
)

type CreateAccountInput struct {
	AccountEmail         string
	DisplayName          string
	AccountSSOEmail      string
	AccountSSOFirstName  string
	AccountSSOLastName   string
	OrganizationalUnitID string
	AccountName          string
}

func CreateAccount(ctx context.Context, accountDetails CreateAccountInput) (string, error) {
	logger := activity.GetLogger(ctx)

	var (
		acceptLanguage         = "en"
		awsRegion              = "us-east-1"
		pathName               = "AWS Control Tower Account Factory Portfolio"
		productID              = "pa-sso2iotgrsvbw"
		productName            = "AWS Control Tower Account Factory"
		provisioningParameters = []types.ProvisioningParameter{
			{
				Key:   aws.String("AccountEmail"),
				Value: &accountDetails.AccountEmail,
			},
			{
				Key:   aws.String("AccountName"),
				Value: &accountDetails.AccountName,
			},
			{
				Key:   aws.String("ManagedOrganizationalUnit"),
				Value: &accountDetails.OrganizationalUnitID,
			},
			{
				Key:   aws.String("SSOUserEmail"),
				Value: &accountDetails.AccountEmail,
			},
			{
				Key:   aws.String("SSOUserFirstName"),
				Value: &accountDetails.AccountEmail,
			},
			{
				Key:   aws.String("SSOUserLastName"),
				Value: &accountDetails.AccountEmail,
			},
		}
	)

	provisionProductInput := servicecatalog.ProvisionProductInput{
		AcceptLanguage:         &acceptLanguage,
		PathName:               &pathName,
		ProductName:            &productName,
		ProvisionedProductName: &accountDetails.AccountName,
		ProvisioningArtifactId: &productID,
		ProvisioningParameters: provisioningParameters,
	}

	logger.Info("Calling AWS Service Catalog with: ", awsutil.Prettify(provisionProductInput))
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if err != nil {
		return "", fmt.Errorf("Unable to load SDK config: %w", err)
	}

	svc := servicecatalog.NewFromConfig(cfg)
	output, err := svc.ProvisionProduct(context.TODO(), &provisionProductInput)
	if err != nil {
		return "", fmt.Errorf("Could not Provision Product: %w", err)
	}
	logger.Info("Service Catalog Responded with: ", awsutil.Prettify(output))
	return awsutil.Prettify(output), nil
}

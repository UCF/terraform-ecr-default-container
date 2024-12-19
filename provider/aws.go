package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func AWSProvider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"aws_ecr_repository": resourceECRRepository(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	//Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %w", err)
	}

	// Create the ECR client
	return ecr.NewFromConfig(cfg), nil
}

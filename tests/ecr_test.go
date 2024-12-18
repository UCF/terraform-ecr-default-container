package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccECRRepository(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { /* Pre-check logic */ },
		Provider: map[string]func() *schema.Provider{
			"aws": func() *schema.Provider {
				return awsProvider()
			},
			"ecr": func() *schema.Provider {
				return ecrProvider()
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "aws" {
						region = "us-east-1"
					}

					resource "aws_ecr_repository" "test" {
						name = "test-repo"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aws_ecr_repository.test", "name", "test-repo"),
				),
			},
		},
	})
}

package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccECRRepository(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() { /* Pre-check logic */ },
		Steps: []resource.TestStep{
			{
				Config: `
					provider "aws" {
						region "us-east-1"
					}

					resource "aws_ecr_repository" "test" {
						name "test-repo"
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("aws_ecr_repository.test", "name", "test-repo"),
				),
			},
		},
	})
}
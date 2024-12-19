package tests

import (
	"testing"

	"github.com/UCF/terraform-ecr-default-container/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccECRRepository(t *testing.T) {
	resourceName := "aws_ecr_repository.test"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { /* Pre-check logic */ },
		Providers: map[string]*schema.Provider{
			"aws": provider.Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccECRRepositoryConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "test-repo"),
				),
			},
		},
	})
}

func testAccECRRepositoryConfig() string {
	return `
	provider "aws" {
		region = "us-east-1"
	}

	resource "aws_ecr_repository" "test" {
		name = "test-repo"
	}`
}

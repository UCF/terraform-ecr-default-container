package tests

import (
	"testing"

	"github.com/UCF/terraform-ecr-default-container/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccECRRepository(t *testing.T) {
	resourceName := "aws_ecr_repository.test"
	repoName := "test-repo"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { /* Pre-check logic */ },
		Providers: map[string]*schema.Provider{
			"aws": provider.AWSProvider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccECRRepositoryConfig(repoName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", repoName),
				),
			},
		},
	})
}

func testAccECRRepositoryConfig(repoName string) string {
	return `
	resource "aws_ecr_repository" "test" {
		name = "` + repoName + `"
	}`
}

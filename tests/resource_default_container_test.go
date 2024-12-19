package tests

import (
	"fmt"
	"testing"

	"github.com/UCF/terraform-ecr-default-container/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDefaultContainer(t *testing.T) {
	resourceName := "default-container.test"
	repoName := "test-repo"
	imageName := "nginx:latest"
	imageName2 := "alpine:latest"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]*schema.Provider{
			"default-container": provider.Provider(),
			"aws":               provider.AWSProvider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDefaultContainerConfig(repoName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "repository_url", fmt.Sprintf("${data.aws_caller_identity.current.account_id}.dkr.ecr.us-east-1.amazonaws.com/%s", repoName)),
					resource.TestCheckResourceAttr(resourceName, "image", imageName),
				),
			},
			{
				Config: testAccDefaultContainerConfig(repoName, imageName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "repository_url", fmt.Sprintf("${data.aws_caller_identity.current.account_id}.dkr.ecr.us-east-1.amazonaws.com/%s", repoName)),
					resource.TestCheckResourceAttr(resourceName, "image", imageName2),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {}

func testAccDefaultContainerConfig(repoName, image string) string {
	return fmt.Sprintf(`
		data "aws_caller_identity" "current" {}

		resource "aws_ecr_repository" "test" {
			name = "%s"
		}

		resource "default-container" "test" {
			repository_url = aws_ecr_repository.test.repository_url
			image           = "%s"
		}`,
		repoName,
		image,
	)
}

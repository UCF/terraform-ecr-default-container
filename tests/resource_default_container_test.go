package tests

import (
	"fmt"
	"regexp"
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

	accountID, err := getAccountID()
	if err != nil {
		t.Fatal("Failed to get AWS account ID:", err)
	}
	if accountID == "" {
		t.Fatal("AWS_ACCOUNT_ID must be set for acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]*schema.Provider{
			"default-container": provider.Provider(),
			"aws":               provider.AWSProvider(),
		},
		Steps: []resource.TestStep{
			{
				Config:  testAccDefaultContainerConfig(repoName, imageName),
				Destroy: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "repository_url", fmt.Sprintf("%s.dkr.ecr.us-east-1.amazonaws.com/%s", accountID, repoName)),
					resource.TestCheckResourceAttr(resourceName, "image", imageName),
				),
			},
			{
				ResourceName:      "default-container.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccDefaultContainerConfig(repoName, imageName2),
				Destroy: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "repository_url", fmt.Sprintf("%s.dkr.ecr.us-east-1.amazonaws.com/%s", accountID, repoName)),
					resource.TestCheckResourceAttr(resourceName, "image", imageName2),
				),
			},
		},
	})
}

func TestAccDefaultContainer_InvalidRepoName(t *testing.T) {
	repoName := "invalid/repo/name"
	imageName := "nginx:latest"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]*schema.Provider{
			"default-container": provider.Provider(),
			"aws":               provider.AWSProvider(),
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccDefaultContainerConfig(repoName, imageName),
				Destroy:     true,
				ExpectError: regexp.MustCompile(`"repository_url" must be a valid repository URL, got: invalid/repo/name`),
			},
		},
	})
}

func TestAccDefaultContainer_InvalidImageName(t *testing.T) {
	repoName := "test-repo"
	imageName := "invalid:image:name"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]*schema.Provider{
			"default-container": provider.Provider(),
			"aws":               provider.AWSProvider(),
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccDefaultContainerConfig(repoName, imageName),
				ExpectError: regexp.MustCompile(`"image" must be a valid image name, got: invalid:image:name`),
			},
		},
	})
}

func TestAccDefaultContainer_ResourceDeletion(t *testing.T) {
	resourceName := "default-container.test"
	repoName := "test-repo"
	imageName := "nginx:latest"

	accountID, err := getAccountID()
	if err != nil {
		t.Fatal("Failed to get AWS account ID:", err)
	}
	if accountID == "" {
		t.Fatal("AWS_ACCOUNT_ID must be set for acceptance tests")
	}

	resource.ParallelTest(
		t,
		resource.TestCase{PreCheck: func() { testAccPreCheck(t) }, Providers: map[string]*schema.Provider{
			"default-container": provider.Provider(),
			"aws":               provider.AWSProvider(),
		}, Steps: []resource.TestStep{
			{
				Config:  testAccDefaultContainerConfig(repoName, imageName),
				Destroy: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						resourceName,
						"repository_url",
						fmt.Sprintf("%s.dkr.ecr.us-east-1.amazonaws.com/%s", accountID, repoName),
					),
					resource.TestCheckResourceAttr(resourceName, "image", imageName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccDefaultContainerConfig(repoName, imageName),
				Destroy: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "repository_url"),
					resource.TestCheckNoResourceAttr(resourceName, "image"),
				),
			},
		}},
	)
}

func testAccPreCheck(t *testing.T) {}

func testAccDefaultContainerConfig(repoName, image string) string {
	return fmt.Sprintf(`
		resource "aws_ecr_repository" "test" {
			name = "%s"
		}

		resource "default-container" "test" {
			repository_url = aws_ecr_repository.test.repository_url
			image           = "%s"

			depends_on = [aws_ecr_repository.test]
		}`,
		repoName,
		image,
	)
}

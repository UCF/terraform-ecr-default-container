package tests

import (
	"fmt"
	"testing"

	"github.com/UCF/terraform-ecr-default-container/provider"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccECRRepository(t *testing.T) {
	resourceName := "aws_ecr_repository.test"
	repoName := "test-repo"

	accountID, err := getAccountID()
	if err != nil {
		t.Fatal("Failed to get AWS account ID:", err)
	}
	if accountID == "" {
		t.Fatal("AWS_ACCOUNT_ID must be set for acceptance tests")
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() { /* Pre-check logic */ },
		Providers: map[string]*schema.Provider{
			"aws": provider.AWSProvider(),
		},
		Steps: []resource.TestStep{
			{
				Config:  testAccECRRepositoryConfig(repoName),
				Destroy: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", repoName),
					resource.TestCheckResourceAttr(resourceName, "repository_url", fmt.Sprintf("%s.dkr.ecr.us-east-1.amazonaws.com/%s", accountID, repoName)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:  testAccECRRepositoryConfig(repoName),
				Destroy: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckNoResourceAttr(resourceName, "name"),
					resource.TestCheckNoResourceAttr(resourceName, "repository_url"),
				),
			},
		},
	})
}

func testAccECRRepositoryConfig(repoName string) string {
	return fmt.Sprintf(`
resource "aws_ecr_repository" "test" {
  name = "%s"
}
`, repoName)
}

func getAccountID() (string, error) {
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		return "", err
	}

	return *result.Account, nil
}

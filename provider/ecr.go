package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceECRRepository() *schema.Resource {
	return &schema.Resource{
		Create: resourceECRRepositoryCreate,
		Read:   resourceECRRepositoryRead,
		Update: resourceECRRepositoryUpdate,
		Delete: resourceECRRepositoryDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !isValidRepoName(v) {
						errs = append(errs, fmt.Errorf("%q must be a valid repository name, got: %s", key, v))
					}
					return
				},
			},
			"repository_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceECRRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*ecr.Client)
	repoName := d.Get("name").(string)

	// Create the ECR repository
	input := &ecr.CreateRepositoryInput{
		RepositoryName: &repoName,
	}
	_, err := client.CreateRepository(context.TODO(), input)
	if err != nil {
		log.Printf("Error creating ECR repository: %v", err)
		return fmt.Errorf("failed to create ECR repository %s: %w", repoName, err)
	}

	d.SetId(repoName)
	return resourceECRRepositoryRead(d, m)
}

func resourceECRRepositoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ecr.Client)
	repoName := d.Id()

	accountID, err := getAccountID()
	if err != nil {
		return fmt.Errorf("Failed to get AWS account ID: %w", err)
	}

	input := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repoName},
	}

	result, err := client.DescribeRepositories(context.TODO(), input)

	if err != nil {
		if isNotFoundError(err) {
			log.Printf("ECR repository %s not found, marking as deleted", repoName)
			d.SetId("")
			return fmt.Errorf("ECR repository is not found")
		}
		return fmt.Errorf("failed to describe ECR repository: %w", err)
	}

	d.Set("name", result.Repositories[0].RepositoryName)

	repoURL := fmt.Sprintf("%s.dkr.ecr.us-east-1.amazonaws.com/%s", accountID, *result.Repositories[0].RepositoryName)
	d.Set("repository_url", repoURL)

	return nil

}

func resourceECRRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceECRRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ecr.Client)
	repoName := d.Id()

	input := &ecr.DeleteRepositoryInput{
		RepositoryName: &repoName,
		Force:          true,
	}

	_, err := client.DeleteRepository(context.TODO(), input)
	return err
}

func isNotFoundError(err error) bool {
	var notFound *types.RepositoryNotFoundException
	return errors.As(err, &notFound)
}

func getAccountID() (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create AWS session: %w", err)
	}
	svc := sts.New(sess)
	input := &sts.GetCallerIdentityInput{}

	result, err := svc.GetCallerIdentity(input)
	if err != nil {
		return "", err
	}

	return *result.Account, nil
}

func isValidRepoName(repoName string) bool {
	// Check length
	if len(repoName) < 2 || len(repoName) > 256 {
		return false
	}

	// Check for valid characters
	var validRepoName = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*[a-z0-9]$`)
	if !validRepoName.MatchString(repoName) {
		return false
	}

	// Check for consecutive periods or hyphens
	if strings.Contains(repoName, "..") || strings.Contains(repoName, "--") || strings.Contains(repoName, "__") {
		return false
	}

	return true
}

package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
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
		return err
	}

	d.SetId(repoName)
	return resourceECRRepositoryRead(d, m)
}

func resourceECRRepositoryRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*ecr.Client)
	repoName := d.Id()

	input := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repoName},
	}

	result, err := client.DescribeRepositories(context.TODO(), input)

	if err != nil {
		if isNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", result.Repositories[0].RepositoryName)

	repoURL := fmt.Sprintf("%s.dkr.ecr.us-east-1.amazonaws.com", *result.Repositories[0].RepositoryName)
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
	return false
}

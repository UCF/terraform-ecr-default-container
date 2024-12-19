package provider

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"default-container":  resourceDefaultContainer(),
			"aws_ecr_repository": resourceECRRepository(),
		},
	}
}

func resourceDefaultContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceDefaultContainerCreate,
		Read:   resourceDefaultContainerRead,
		Update: resourceDefaultContainerUpdate,
		Delete: resourceDefaultContainerDelete,
		Schema: map[string]*schema.Schema{
			"repository_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDefaultContainerCreate(d *schema.ResourceData, m interface{}) error {
	image := d.Get("image").(string)
	repoURL := d.Get("repository_url").(string)

	// Extract the repository name from the repository URL
	// Assuming the URL format is "<repo-name>.dkr.ecr.<region>.amazonaws.com"
	parts := strings.Split(repoURL, ".")
	if len(parts) < 3 {
		return fmt.Errorf("invalid repository URL format: %s", repoURL)
	}
	repoName := parts[0] // This is the repository name

	// Set the repository_name in the resource data
	d.Set("repository_name", repoName)

	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)
	}

	svc := sts.New(sess)
	_, err = svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("failed to get caller identity: %v", err)
	}

	region := "us-east-1"

	authCommand := fmt.Sprintf("aws ecr get-login-password --region %s | podman login --username AWS --password-stdin %s", region, repoURL)

	cmd := exec.Command("bash", "-c", authCommand)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to authenticate Podman with ECR: %v, output: %s", err, output)
	}

	pullCommand := fmt.Sprintf("podman pull %s", image)
	if err := exec.Command("bash", "-c", pullCommand).Run(); err != nil {
		return fmt.Errorf("failed to pull image: %v", err)
	}

	d.SetId(image)
	return nil
}

func resourceDefaultContainerRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDefaultContainerUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceDefaultContainerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*ecr.Client)

	repoURL := d.Get("repository_url").(string)
	image := d.Get("image").(string)

	parts := strings.Split(repoURL, ".")
	if len(parts) < 3 {
		return fmt.Errorf("invalid repository URL format: %s", repoURL)
	}
	repoName := parts[0]

	input := &ecr.BatchDeleteImageInput{
		RepositoryName: &repoName,
		ImageIds: []types.ImageIdentifier{
			{
				ImageTag: aws.String(image),
			},
		},
	}

	_, err := client.BatchDeleteImage(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("failed to delete image from ECR: %v", err)
	}

	d.SetId("")
	return nil
}

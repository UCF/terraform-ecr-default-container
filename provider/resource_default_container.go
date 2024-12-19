package provider

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
	"github.com/aws/aws-sdk-go/aws"
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
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !isValidRepoURL(v) {
						errs = append(errs, fmt.Errorf("%q must be a valid repository URL, got: %s", key, v))
					}
					return
				},
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !isValidImageName(v) {
						errs = append(errs, fmt.Errorf("%q must be a valid image name, got: %s", key, v))
					}
					return
				},
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
	parts := strings.Split(repoURL, "/")
	repoName := parts[len(parts)-1]

	// Set the repository_name in the resource data
	d.Set("repository_name", repoName)

	region := "us-east-1"

	authCommand := fmt.Sprintf("aws ecr get-login-password --region %s | podman login --username AWS --password-stdin %s", region, repoURL)

	cmd := exec.Command("bash", "-c", authCommand)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("failed to authenticate Podman with ECR: %v, output: %s", err, output)
	}

	pullCommand := fmt.Sprintf("podman pull %s", image)
	if err := exec.Command("bash", "-c", pullCommand).Run(); err != nil {
		return fmt.Errorf("failed to pull image with command: %s %v", pullCommand, err)
	}

	pushCommand := fmt.Sprintf("podman push %s", repoURL)
	if err := exec.Command("bash", "-c", pushCommand).Run(); err != nil {
		return fmt.Errorf("failed to push image to ECR with command: %s %v", pushCommand, err)
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

	parts := strings.Split(repoURL, "/")
	repoName := parts[len(parts)-1]

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

func isValidImageName(imageName string) bool {
	// Simple regex to validate image name format
	var validImageName = regexp.MustCompile(`^[a-zA-Z0-9._-]+(:[a-zA-Z0-9._-]+)?$`)
	return validImageName.MatchString(imageName)
}

func isValidRepoURL(repoURL string) bool {
	// Simple regex to validate repository URL format
	var validRepoURL = regexp.MustCompile(`^[a-zA-Z0-9._-]+\.dkr\.ecr\.[a-zA-Z0-9-]+\.amazonaws\.com/[a-zA-Z0-9._-]+$`)
	return validRepoURL.MatchString(repoURL)
}

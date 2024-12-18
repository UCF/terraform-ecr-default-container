package provider

import (
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
		},
	}
}

func resourceECRRepositoryCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceECRRepositoryRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceECRRepositoryUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceECRRepositoryDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws"
)

func Provider() *schema.Provider {
	return aws.Provider()
}

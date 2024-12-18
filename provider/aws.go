package provider

import (
	"github.com:hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com:hashicorp/terraform-provider-aws/aws"
)

func awsProvider() *schema.Provider {
	return aws.Provider()
}

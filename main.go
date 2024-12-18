package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		Terraform: terraform.ProviderFunc(Provider),
	})
}

func Provider() *schema.Provider {
}

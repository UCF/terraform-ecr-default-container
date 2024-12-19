package main

import (
	"github.com/UCF/terraform-ecr-default-container/provider"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
		GRPCProviderFunc: func() tfprotov5.ProviderServer {
		},
		GRPCProviderV6Func: func() tfprotov6.ProviderServer {
		},
		Logger:              nil,
		Debug:               false,
		NoLogOutputOverride: false,
		UseTFLogSink:        nil,
		ProviderAddr:        "",
	})
}

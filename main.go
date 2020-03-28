package main

import (
	"github.com/b9company/terraform-provider-s3/s3"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: s3.Provider,
	})
}

package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/lexfrei/terraform-provider-windns/windns"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: windns.Provider,
	})
}

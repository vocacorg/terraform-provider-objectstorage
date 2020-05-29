package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/vocacorg/terraform-provider-objectstorage/objectstorage"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: objectstorage.Provider})
}

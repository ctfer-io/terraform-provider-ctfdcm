package provider_test

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	tfctfd "github.com/ctfer-io/terraform-provider-ctfd/provider"
	"github.com/ctfer-io/terraform-provider-ctfdcm/provider"
)

const (
	providerConfig = `
provider "ctfd" {}
provider "ctfdcm" {}
`
)

var (
	// testAccProtoV6ProviderFactories are used to instantiate a provider during
	// acceptance testing. The factory function will be invoked for every Terraform
	// CLI command executed to create a provider server to which the CLI can
	// reattach.
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"ctfdcm": providerserver.NewProtocol6WithError(provider.New("test")()),
		"ctfd":   providerserver.NewProtocol6WithError(tfctfd.New("test")()),
	}
)

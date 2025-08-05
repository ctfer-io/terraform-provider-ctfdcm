package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_Instance_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "ctfdcm_challenge_dynamiciac" "chall" {
	name        = "Some challenge"
	category    = "cat"
	description = "..."
	value       = 500
    decay       = 20
    minimum     = 50
    state       = "visible"

	scenario = var.scenario
}

resource "ctfd_user" "pandatix" {
	name     = "PandatiX"
	email    = "lucastesson@protonmail.com"
	password = "password"
}

resource "ctfd_team" "ctfer" {
	name = "CTFer.io"
	email = "ctfer-io@protonmail.com"
	password = "ctfer"
	members = [
	  ctfd_user.pandatix.id,
	]
	captain = ctfd_user.pandatix.id
}

resource "ctfdcm_instance" "ist" {
	challenge_id = ctfdcm_challenge_dynamiciac.chall.id
	source_id = ctfd_team.ctfer.id
}

variable "scenario" {
  type = string
}
`,
				ConfigVariables: config.Variables{
					"scenario": config.StringVariable(ref),
				},
			},
		},
	})
}

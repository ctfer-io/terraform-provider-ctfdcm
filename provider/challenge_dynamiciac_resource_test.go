package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAcc_ChallengeDynamicIaC_Lifecycle(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "ctfdcm_challenge_dynamiciac" "http" {
	name        = "HTTP Authentication"
	category    = "network"
	description = <<-EOT
        Oh no ! I did not see my connection was no encrypted !
        I hope no one spied me...
    EOT
	attribution = "Nicolas"
	value       = 500
    decay       = 20
    minimum     = 50
    state       = "hidden"

	scenario = var.scenario

	topics = [
		"Network"
	]
	tags = [
		"network"
	]
}

variable "scenario" {
  type = string
}
`,
				ConfigVariables: config.Variables{
					"scenario": config.StringVariable(ref),
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ctfdcm_challenge_dynamiciac.http", "id"),
					resource.TestCheckNoResourceAttr("ctfdcm_challenge_dynamiciac.http", "timeout"),
					resource.TestCheckNoResourceAttr("ctfdcm_challenge_dynamiciac.http", "until"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ctfdcm_challenge_dynamiciac.http",
				ImportState:       true,
				ImportStateVerify: true,
				ConfigVariables: config.Variables{
					"scenario": config.StringVariable(ref),
				},
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "ctfdcm_challenge_dynamiciac" "http" {
	name        = "HTTP Authentication"
	category    = "network"
	description = <<-EOT
        Oh no ! I did not see my connection was no encrypted !
        I hope no one spied me...
    EOT
	attribution = "Nicolas"
	value       = 500
    decay       = 20
    minimum     = 50
    state       = "hidden"

	shared          = true
    destroy_on_flag = true
    mana_cost       = 1
    scenario        = var.scenario
    timeout         = 600
	additional      = {
	    key = "value"
	}

	min = 2
	max = 4

	topics = [
		"Network"
	]
	tags = [
		"network"
	]
}

variable "scenario" {
  type = string
}
`,
				ConfigVariables: config.Variables{
					"scenario": config.StringVariable(ref),
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

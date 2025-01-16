package provider_test

import (
	"testing"

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

	scenario_id = ctfd_file.scenario.id

	topics = [
		"Network"
	]
	tags = [
		"network"
	]
}

resource "ctfd_file" "scenario" {
  name         = "scenario.zip"
  contentb64   = filebase64("some zip content")
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ctfdcm_challenge_dynamiciac.http", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "ctfdcm_challenge_dynamic.http",
				ImportState:       true,
				ImportStateVerify: true,
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
    scenario_id     = ctfd_file.scenario.id
    timeout         = 600

	topics = [
		"Network"
	]
	tags = [
		"network"
	]
}

resource "ctfd_file" "scenario" {
  name         = "scenario.zip"
  contentb64   = filebase64("some zip content")
}
`,
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

package provider_test

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"gopkg.in/yaml.v2"
)

func TestAcc_ChallengeDynamicIaC_Lifecycle(t *testing.T) {
	scn := config.StringVariable(scenario())

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
  name       = "scenario.zip"
  contentb64 = var.scenario
}

variable "scenario" {
  type = string
}
`,
				ConfigVariables: config.Variables{
					"scenario": scn,
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("ctfdcm_challenge_dynamiciac.http", "id"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "ctfdcm_challenge_dynamiciac.http",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"ctfd_file.scenario"},
				ConfigVariables: config.Variables{
					"scenario": scn,
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
    scenario_id     = ctfd_file.scenario.id
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

resource "ctfd_file" "scenario" {
  name       = "scenario.zip"
  contentb64 = var.scenario
}

variable "scenario" {
  type = string
}
`,
				ConfigVariables: config.Variables{
					"scenario": scn,
				},
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func scenario() string {
	buf := bytes.NewBuffer([]byte{})
	archive := zip.NewWriter(buf)

	// Add Pulumi.yaml file
	mp := map[string]any{
		"name": "scenario",
		"runtime": map[string]any{
			"name": "go",
			"options": map[string]any{
				"binary": "./main",
			},
		},
		"description": "An example scenario.",
	}
	b, err := yaml.Marshal(mp)
	if err != nil {
		panic(err)
	}
	w, err := archive.Create("Pulumi.yaml")
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(w, bytes.NewBuffer(b)); err != nil {
		panic(err)
	}

	// Add binary file
	fs, err := compile()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = fs.Close()
	}()

	fst, err := fs.Stat()
	if err != nil {
		panic(err)
	}
	header, err := zip.FileInfoHeader(fst)
	if err != nil {
		panic(err)
	}
	header.Name = "main"

	// Create archive
	f, err := archive.CreateHeader(header)
	if err != nil {
		panic(err)
	}

	// Copy the file's contents into the archive.
	_, err = io.Copy(f, fs)
	if err != nil {
		panic(err)
	}

	// Complete zip creation
	if err := archive.Close(); err != nil {
		panic(err)
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes())
}

func compile() (*os.File, error) {
	cmd := exec.Command("go", "build", "-o", "main", "../scenario/main.go")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	defer func() {
		cmd := exec.Command("rm", "main")
		_ = cmd.Run()
	}()
	return os.Open("main")
}

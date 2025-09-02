package provider_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ctfer-io/chall-manager/pkg/scenario"
)

var (
	REGISTRY = ""

	ref = "scenario:v0.1.0"
)

func TestMain(m *testing.M) {
	r, ok := os.LookupEnv("REGISTRY")
	if !ok {
		fmt.Println("Environment variable REGISTRY is not set, please indicate the domain name/IP address to reach out the registry.")
		os.Exit(1)
	}
	REGISTRY = r
	ref = fmt.Sprintf("%s/%s", REGISTRY, ref)

	if err := pushScenario(); err != nil {
		fmt.Printf("Pushing scenario %s: %s", ref, err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func pushScenario() error {
	ctx := context.Background()
	return scenario.EncodeOCI(ctx, ref, "./scenario", true, "", "")
}

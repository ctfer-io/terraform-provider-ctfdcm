package provider_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ctfer-io/chall-manager/pkg/scenario"
	tfctfd "github.com/ctfer-io/terraform-provider-ctfd/v2/provider"
	"github.com/ctfer-io/terraform-provider-ctfdcm/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	REGISTRY = ""

	ref = "scenario:v0.1.0"
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// Prepare OTel tracing
	out, err := provider.SetupOTelSDK(ctx, "test")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := out.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	testAccProtoV6ProviderFactories["ctfd"] = providerserver.NewProtocol6WithError(tfctfd.New("test", out.TracerProvider)())
	testAccProtoV6ProviderFactories["ctfdcm"] = providerserver.NewProtocol6WithError(provider.New("test", out.TracerProvider)())

	// Build and push test scenario
	r, ok := os.LookupEnv("REGISTRY")
	if !ok {
		fmt.Println("Environment variable REGISTRY is not set, please indicate the domain name/IP address to reach out the registry.")
		os.Exit(1)
	}
	REGISTRY = r
	ref = fmt.Sprintf("%s/%s", REGISTRY, ref)

	if err := func() error {
		ctx, span := out.TracerProvider.Tracer("terraform-provider-ctfdcm").Start(ctx, "push-scenario")
		defer span.End()

		return scenario.EncodeOCI(ctx, ref, "./scenario", true, "", "")
	}(); err != nil {
		fmt.Printf("Pushing scenario %s: %s", ref, err)
		os.Exit(1)
	}

	if sc := m.Run(); sc != 0 {
		log.Fatalf("Failed with status code %d", sc)
	}
}

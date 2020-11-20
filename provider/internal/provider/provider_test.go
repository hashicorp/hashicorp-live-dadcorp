package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	tfmux "github.com/hashicorp/terraform-plugin-mux"
)

var testProviders = map[string]func() (tfprotov5.ProviderServer, error){
	"dadcorp": func() (tfprotov5.ProviderServer, error) {
		ctx := context.Background()
		sdkv2 := New().GRPCProvider
		plugin := NewPlugin

		factory, err := tfmux.NewSchemaServerFactory(ctx, sdkv2, plugin)
		if err != nil {
			return nil, err
		}
		return factory.Server(), nil
	},
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("DADCORP_USERNAME") == "" {
		t.Fatalf("DADCORP_USERNAME must be set")
	}
	if os.Getenv("DADCORP_PASSWORD") == "" {
		t.Fatalf("DADCORP_PASSWORD must be set")
	}
}

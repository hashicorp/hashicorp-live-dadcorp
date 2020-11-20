package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testProviders = map[string]func() (*schema.Provider, error){
	"dadcorp": func() (*schema.Provider, error) {
		return New(), nil
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

package main

import (
	"context"
	"log"
	"os"

	"dadcorp.dev/terraform-provider-dadcorp/internal/provider"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	tf5server "github.com/hashicorp/terraform-plugin-go/tfprotov5/server"
	tfmux "github.com/hashicorp/terraform-plugin-mux"
)

func main() {
	ctx := context.Background()
	sdkv2 := provider.New().GRPCProvider

	factory, err := tfmux.NewSchemaServerFactory(ctx, sdkv2)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
	tf5server.Serve("registry.terraform.io/dadcorp/dadcorp", func() tfprotov5.ProviderServer {
		return factory.Server()
	})
}

module dadcorp.dev/terraform-provider-dadcorp

go 1.15

require (
	dadcorp.dev/client v0.0.0-00010101000000-000000000000
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.2.0
)

replace dadcorp.dev/client => ../client

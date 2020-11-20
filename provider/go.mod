module dadcorp.dev/terraform-provider-dadcorp

go 1.15

require (
	dadcorp.dev/client v0.0.0-00010101000000-000000000000
	github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.2.0
	github.com/nsf/jsondiff v0.0.0-20200515183724-f29ed568f4ce
)

replace dadcorp.dev/client => ../client

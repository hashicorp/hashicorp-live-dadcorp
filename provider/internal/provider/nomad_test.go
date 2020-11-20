package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNomadCluster_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testProviders,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigNomadCluster_basic(),
			},
			{
				ResourceName:      "dadcorp_nomad_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigNomadCluster_updated(),
			},
			{
				ResourceName:      "dadcorp_nomad_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigNomadCluster_basic() string {
	return `
resource "dadcorp_nomad_cluster" "test" {
  name = "test cluster"
  datacenter = "dc1"

  advertise {
    http = "0.0.0.0"
    rpc = "0.0.0.0"
    serf = "0.0.0.0"
  }

  ports {
    http = 1234
    rpc = 2345
    serf = 3456
  }

  server {
    server_join {
      retry_join = ["1.2.3.4"]
      start_join = ["2.3.4.6"]
      retry_max = 5
      retry_interval = "15s"
    }
  }
}
`
}

func testAccConfigNomadCluster_updated() string {
	return `
resource "dadcorp_nomad_cluster" "test" {
  name = "test cluster updated"
  datacenter = "dc2"
  bind_addr = "1.0.1.0"

  advertise{}

  ports {
    http = 2345
    rpc = 3456
    serf = 4567
  }

  server {
    server_join {
      retry_join = ["1.2.3.4", "1.1.1.1"]
      start_join = ["2.3.4.5", "2.2.2.2"]
      retry_max = 2
      retry_interval = "1m"
    }
  }
}
`
}

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccConsulCluster_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigConsulCluster_basic(),
			},
			{
				ResourceName:      "dadcorp_consul_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigConsulCluster_updated(),
			},
			{
				ResourceName:      "dadcorp_consul_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigConsulCluster_basic() string {
	return `
resource "dadcorp_consul_cluster" "test" {
  name = "test cluster"
  bind_addr = "1.2.3.4"
  addresses {
    dns = "127.0.0.1"
    http = "0.0.0.0"
    https = "0.0.0.0"
    grpc = "0.0.0.0"
  }

  ports {
    dns = 18600
    http = 18500
    https = 18501
    grpc = 18502
    serf_lan = 18301
    serf_wan = 18302
    server = 18300
    sidecar_min_port = 20000
    sidecar_max_port = 21000
    expose_min_port = 30000
    expose_max_port = 31000
  }
}
`
}

func testAccConfigConsulCluster_updated() string {
	return `
resource "dadcorp_consul_cluster" "test" {
  name = "test cluster"
  bind_addr = "2.3.4.5"
  addresses {
    dns = "1.2.3.4"
    http = "5.6.7.8"
    https = "9.10.11.12"
    grpc = "13.14.15.16"
  }

  ports {
    dns = 1
    http = 2
    https = 3
    grpc = 4
    serf_lan = 5
    serf_wan = 6
    server = 7
    sidecar_min_port = 8
    sidecar_max_port = 9
    expose_min_port = 10
    expose_max_port = 11
  }
}
`
}

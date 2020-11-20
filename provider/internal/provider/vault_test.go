package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccVaultCluster_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testProviders,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigVaultCluster_basic(),
			},
			{
				ResourceName:      "dadcorp_vault_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigVaultCluster_updated(),
			},
			{
				ResourceName:      "dadcorp_vault_cluster.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigVaultCluster_basic() string {
	return `
resource "dadcorp_vault_cluster" "test" {
  name = "test cluster"
  region = "us-va-1"
  default_lease_ttl = "1h"
  max_lease_ttl = "24h"

  tcp_listener {
    address = "1.2.3.4"
    cluster_address = "2.3.4.5"
  }
}
`
}

func testAccConfigVaultCluster_updated() string {
	return `
resource "dadcorp_vault_cluster" "test" {
  name = "updated test cluster"
  region = "us-va-2"
  default_lease_ttl = "2h"
  max_lease_ttl = "48h"

  tcp_listener {
    address = "2.3.4.5"
    cluster_address = "3.4.5.6"
  }
}
`
}

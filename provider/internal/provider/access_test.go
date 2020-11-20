package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAccessPolicy_consul(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAccessPolicy_consul_basic(),
			},
			{
				ResourceName:      "dadcorp_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigAccessPolicy_consul_updated(),
			},
			{
				ResourceName:      "dadcorp_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigAccessPolicy_consul_basic() string {
	return `
data "dadcorp_access_policy_consul" "test" {
  cluster_id = "test"
  key = "demo"
  read = true
  write = true
  delete = true
}

resource "dadcorp_access_policy" "test" {
  type = "consul"
  policy_data = data.dadcorp_access_policy_consul.test.json
}
`
}

func testAccConfigAccessPolicy_consul_updated() string {
	return `
data "dadcorp_access_policy_consul" "test" {
  cluster_id = "test"
  key = "demo"
  read = false
  write = false
  delete = false
}

resource "dadcorp_access_policy" "test" {
  type = "consul"
  policy_data = data.dadcorp_access_policy_consul.test.json
}
`
}

func TestAccAccessPolicy_nomad(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAccessPolicy_nomad_basic(),
			},
			{
				ResourceName:      "dadcorp_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigAccessPolicy_nomad_updated(),
			},
			{
				ResourceName:      "dadcorp_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigAccessPolicy_nomad_basic() string {
	return `
data "dadcorp_access_policy_nomad" "test" {
  cluster_id = "test"
  submit_jobs = true
  read_job_status = true
  cancel_jobs = true
}

resource "dadcorp_access_policy" "test" {
  type = "nomad"
  policy_data = data.dadcorp_access_policy_nomad.test.json
}
`
}

func testAccConfigAccessPolicy_nomad_updated() string {
	return `
data "dadcorp_access_policy_nomad" "test" {
  cluster_id = "test"
  submit_jobs = false
  read_job_status = false
  cancel_jobs = false
}

resource "dadcorp_access_policy" "test" {
  type = "nomad"
  policy_data = data.dadcorp_access_policy_nomad.test.json
}
`
}

func TestAccAccessPolicy_terraform(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAccessPolicy_terraform_basic(),
			},
			{
				ResourceName:      "dadcorp_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigAccessPolicy_terraform_updated(),
			},
			{
				ResourceName:      "dadcorp_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigAccessPolicy_terraform_basic() string {
	return `
data "dadcorp_access_policy_terraform" "test" {
  workspace_id = "test"
  plan = true
  apply = true
  override_policies = true
}

resource "dadcorp_access_policy" "test" {
  type = "terraform"
  policy_data = data.dadcorp_access_policy_terraform.test.json
}
`
}

func testAccConfigAccessPolicy_terraform_updated() string {
	return `
data "dadcorp_access_policy_terraform" "test" {
  workspace_id = "test"
  plan = false
  apply = false
  override_policies = false
}

resource "dadcorp_access_policy" "test" {
  type = "terraform"
  policy_data = data.dadcorp_access_policy_terraform.test.json
}
`
}

func TestAccAccessPolicy_vault(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProviderFactories: testProviders,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigAccessPolicy_vault_basic(),
			},
			{
				ResourceName:      "dadcorp_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigAccessPolicy_vault_updated(),
			},
			{
				ResourceName:      "dadcorp_access_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigAccessPolicy_vault_basic() string {
	return `
data "dadcorp_access_policy_vault" "test" {
  cluster_id = "test"
  key = "demo"
  read = true
  write = true
  delete = true
}

resource "dadcorp_access_policy" "test" {
  type = "vault"
  policy_data = data.dadcorp_access_policy_vault.test.json
}
`
}

func testAccConfigAccessPolicy_vault_updated() string {
	return `
data "dadcorp_access_policy_vault" "test" {
  cluster_id = "test"
  key = "demo"
  read = false
  write = false
  delete = false
}

resource "dadcorp_access_policy" "test" {
  type = "vault"
  policy_data = data.dadcorp_access_policy_vault.test.json
}
`
}

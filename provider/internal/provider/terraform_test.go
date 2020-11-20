package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTerraformWorkspace_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testProviders,
		PreCheck:                 func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTerraformWorkspace_basic(),
			},
			{
				ResourceName:      "dadcorp_terraform_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTerraformWorkspace_updated(),
			},
			{
				ResourceName:      "dadcorp_terraform_workspace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigTerraformWorkspace_basic() string {
	return `
resource "dadcorp_terraform_workspace" "test" {
  name = "test workspace"
  agent_pool_id = "test-pool"
  allow_destroy_plan = true
  auto_apply = false
  description = "test workspace from Terraform"
  execution_mode = "agent"
  file_triggers_enabled = true
  queue_all_runs = false
  speculative_enabled = true
  terraform_version = "0.14.0-beta1"
  trigger_prefixes = ["test", "test2"]
  working_directory = "my_dir"

  vcs_repo {
    oauth_token_id = "test-token"
    branch = "master"
    ingress_submodules = true
    identifier = "hashicorp/test"
  }
}
`
}

func testAccConfigTerraformWorkspace_updated() string {
	return `
resource "dadcorp_terraform_workspace" "test" {
  name = "test workspace updated"
  allow_destroy_plan = false
  auto_apply = true
  description = "updated workspace from Terraform"
  execution_mode = "remote"
  file_triggers_enabled = false
  queue_all_runs = true
  speculative_enabled = false
  terraform_version = "0.13.5"
  trigger_prefixes = ["test0", "test2", "test3"]
  working_directory = "my_other_dir"

  vcs_repo {
    oauth_token_id = "test-token-2"
    branch = "main"
    ingress_submodules = false
    identifier = "hashicorp/test-repo"
  }
}
`
}

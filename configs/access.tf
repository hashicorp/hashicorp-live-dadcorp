data "dadcorp_access_policy_terraform" "demo" {
  workspace_id      = "hashicorp-live"
  plan              = true
  apply             = true
  override_policies = false
}

resource "dadcorp_access_policy" "demo" {
  type        = "terraform"
  policy_data = data.dadcorp_access_policy_terraform.demo.json
}

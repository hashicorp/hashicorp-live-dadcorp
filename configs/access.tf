resource "dadcorp_access_policy" "demo" {
  type        = "terraform"
  policy_data = {
    workspace_id = "hashicorp-live"
    plan = true
    apply = true
    override_policies = false
  }
}

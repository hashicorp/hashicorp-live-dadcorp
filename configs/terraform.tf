resource "dadcorp_terraform_workspace" "demo" {
  name             = "hashicorp-live"
  agent_pool_id    = "hashicorp-live-agent-pool"
  execution_mode   = "agent"
  trigger_prefixes = ["config", "examples", "legacy"]

  vcs_repo {
    oauth_token_id = "hashicorp-live"
    branch         = "main"
    identifier     = "dadcorp/demo"
  }
}

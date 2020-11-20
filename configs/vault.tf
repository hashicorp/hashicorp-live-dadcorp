resource "dadcorp_vault_cluster" "demo" {
  name   = "hashicorp-live"
  region = "us-va-2"
  tcp_listener {}
}

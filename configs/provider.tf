provider "dadcorp" {
  username = "admin"
  password = "hunter2"
}

terraform {
  required_providers {
    dadcorp = {
      versions = ["0.1.0"]
      source   = "registry.terraform.io/dadcorp/dadcorp"
    }
  }
}

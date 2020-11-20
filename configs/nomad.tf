#resource "dadcorp_ip" "demo" {}

resource "dadcorp_nomad_cluster" "demo" {
  name       = "hashicorp-live"
  datacenter = "dc2"
  #	bind_addr  = dadcorp_ip.demo.ip

  advertise {
    http = "0.1.0.1"
    rpc  = "2.0.0.2"
    serf = "0.1.2.0"
  }

  ports {
    http = 2345
    rpc  = 3456
    serf = 4567
  }

  server {
    server_join {
      retry_join     = ["1.2.3.4", "1.1.1.1"]
      start_join     = ["2.3.4.5", "2.2.2.2"]
      retry_max      = 2
      retry_interval = "1m"
    }
  }
}

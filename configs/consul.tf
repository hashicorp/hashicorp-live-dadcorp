resource "dadcorp_consul_cluster" "demo" {
  name = "hashicorp-live"
  # bind_addr = "1.2.3.4"
  addresses {
    dns   = "127.0.0.1"
    http  = "0.0.0.0"
    https = "0.0.0.0"
    grpc  = "0.0.0.0"
  }

  ports {
    dns              = 18600
    http             = 18500
    https            = 18501
    grpc             = 18502
    serf_lan         = 18301
    serf_wan         = 18302
    server           = 18300
    sidecar_min_port = 20000
    sidecar_max_port = 21000
    expose_min_port  = 30000
    expose_max_port  = 31000
  }
}

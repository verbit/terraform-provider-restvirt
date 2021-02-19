resource "restvirt_domain" "test" {
  name       = "my-vm"
  vcpu       = 4
  memory     = 2048
  private_ip = "192.168.123.168"
  user_data  = file("${path.module}/cloud-config.yaml")
}

resource "restvirt_forwarding" "test" {
  source_port = 2222
  target_ip   = restvirt_domain.test.private_ip
  target_port = 22
}

resource "restvirt_domain" "test" {
  name       = "my-vm"
  vcpu       = 4
  memory     = 2048
  private_ip = "192.168.123.168"
  user_data  = file("${path.module}/cloud-config.yaml")
}

resource "restvirt_volume" "test" {
  name = "vol-test"
  size = 10
}

resource "restvirt_attachment" "test" {
  domain_id = restvirt_domain.test.id
  volume_id = restvirt_volume.test.id
}

output "disk_address" {
  value = restvirt_attachment.test.disk_address
}

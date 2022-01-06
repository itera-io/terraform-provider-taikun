resource "taikun_standalone_profile" "foo" {
  name       = "foo"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"

  security_group {
    name        = "http"
    from_port   = 80
    to_port     = 80
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
  security_group {
    name        = "https"
    from_port   = 443
    to_port     = 443
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
  security_group {
    name        = "dhcp"
    from_port   = 67
    to_port     = 68
    ip_protocol = "udp"
    cidr        = "0.0.0.0/0"
  }
  security_group {
    name        = "icmp"
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }

  lock            = true
  organization_id = 42
}
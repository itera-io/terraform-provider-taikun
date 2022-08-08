resource "taikun_access_profile" "foo" {
  name = "foo"

  organization_id = "42"
  lock            = true
  http_proxy      = "proxy_url"

  ssh_user {
    name       = "oui oui"
    public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGQwGpzLk0IzqKnBpaHqecLA+X4zfHamNe9Rg3CoaXHF :oui_oui:"
  }

  ntp_server {
    address = "time.windows.com"
  }

  dns_server {
    address = "8.8.8.8"
  }

  dns_server {
    address = "8.8.4.4"
  }

  allowed_host {
    description = "Host A"
    address     = "10.0.0.1"
    mask_bits   = 24
  }
}

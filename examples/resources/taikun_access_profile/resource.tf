resource "taikun_access_profile" "foo" {
  # Required
  name = "foo"

  # Optional
  organization_id = "42" # Optional for Partner and Admin
  is_locked       = true
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
}
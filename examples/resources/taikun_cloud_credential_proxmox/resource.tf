resource "taikun_cloud_credential_proxmox" "foo" {
  name = "foo"

  api_host         = "https://foo.example.com/fooapi/json"
  client_id        = "foo@bar.example.com"
  client_secret    = "longAND420&#^complicatedPASSWORD"
  hypervisors      = ["foo-proxmox", "foo-proxmox2"]
  storage          = "lvm-thin"
  vm_template_name = "foo-template"

  private_ip_address             = "192.168.0.0"
  private_net_mask               = "24"
  private_gateway                = "192.168.0.1"
  private_begin_allocation_range = "192.168.0.10"
  private_end_allocation_range   = "192.168.0.20"
  private_bridge                 = "brvm1"

  public_ip_address             = "66.55.44.0"
  public_net_mask               = "24"
  public_gateway                = "66.55.44.1"
  public_begin_allocation_range = "66.55.44.33"
  public_end_allocation_range   = "66.55.44.44"
  public_bridge                 = "brvm2"

  organization_id = "42"
  lock            = false
  continent       = "Europe"
}

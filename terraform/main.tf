terraform {
  required_providers {
    ibm = {
      source = "IBM-Cloud/ibm"
      version = ">= 1.12.0"
    }
  }
}

provider "ibm" {
  ibmcloud_api_key      = var.ibmcloud_api_key
  iaas_classic_username = var.iaas_classic_username
  iaas_classic_api_key  = var.iaas_classic_api_key
}

resource "ibm_network_vlan" "ci_vlans" {
  count = var.vlan_quantity

  name            = "ci_vlan_${count.index}"
  datacenter      = var.datacenter
  router_hostname = var.router
  type            = "PRIVATE"
  tags = var.vlan_tags
}


resource "ibm_subnet" "portable_subnet" {
  count = var.vlan_quantity

  type       = "Portable"
  private    = true
  ip_version = 4
  capacity   = var.subnet_capacity
  vlan_id    = ibm_network_vlan.ci_vlans[count.index].id
  notes      = "portable_subnet"
  tags = var.vlan_tags

  //User can increase timeouts
  timeouts {
    create = "45m"
  }
}

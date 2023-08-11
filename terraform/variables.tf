variable "ibmcloud_api_key" {
  type      = string
  sensitive = true
}

variable "iaas_classic_username" {
  type      = string
  sensitive = true
}

variable "iaas_classic_api_key" {
  type      = string
  sensitive = true
}

variable "vlan_quantity" {
  type = number
  default = 1
}

variable "subnet_capacity" {
  type = number
  default = 128
}

variable "datacenter" {
  type = string
}

variable "router" {
  type = string
}

variable "vlan_tags" {
  type = list(string)
}

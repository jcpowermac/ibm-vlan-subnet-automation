package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

type Vlan struct {
	VlanId int `json:"vlanId"`
	Subnet `json:"subnet"`
}

type Subnet struct {
	Cidr               int      `json:"cidr"`
	DnsServer          string   `json:"dnsServer"`
	MachineNetworkCidr string   `json:"machineNetworkCidr"`
	Gateway            string   `json:"gateway"`
	Mask               string   `json:"mask"`
	Network            string   `json:"network"`
	IpAddresses        []string `json:"ipAddresses"`
}

func main() {

	/*
	 * $notes = "gateway: %s, cidr: %d, vlan: %d"
	 * New-VDPortgroup -Name "ci-vlan-%d" -Notes $notes -VDSwitch $vdswitch -VLanId %d
	 */

	powercli := `
$notes = "vlan: {{.VlanId}} gateway: {{.Subnet.Gateway}} cidr: {{.Subnet.Cidr}} mask: {{.Subnet.Mask}}"
New-VDPortgroup -Name "ci-vlan-{{.VlanId}}" -Notes $notes -VDSwitch $vdswitch -VLanId {{.VlanId}} 
`

	interfaces := `
	   	 interfaces {
	              bonding dp0bond0 {
	                      vif {{.VlanId }} {
	                              address 192.168.93.1/30
	                              vrrp {
	                                      vrrp-group 1 {
	                                              preempt false
	                                              priority 254
	                                              sync-group vgroup1
	                                              virtual-address {{.Subnet.Gateway}}/{{.Subnet.Cidr}}
	                                      }
	                              }
	                      }
	              }
`

	pgTemplate, err := template.New("pg").Parse(powercli)

	if err != nil {
		log.Fatal(err)
	}
	intTemplate, err := template.New("interface").Parse(interfaces)

	if err != nil {
		log.Fatal(err)
	}

	subnetvlanmap := make(map[int]Subnet)

	// https://github.com/softlayer/softlayer-go/blob/master/examples/cmd/vlan_detail.go
	// https://sldn.softlayer.com/reference/datatypes/SoftLayer_Network_Subnet/

	objectMask :=
		`mask[id,name,vlanNumber,subnets[id,cidr,netmask,networkIdentifier,subnetType]]`
	subnetObjectMask :=
		`mask[id,ipAddressCount,gateway,cidr,netmask,networkIdentifier,ipAddresses[ipAddress,isNetwork,isBroadcast,isGateway]]`

	sess := session.New()

	service := services.GetAccountService(sess)
	networkVlanService := services.GetNetworkVlanService(sess)

	vlans, err := service.GetNetworkVlans()

	if err != nil {
		log.Fatal(err)
	}

	networkSubnetService := services.GetNetworkSubnetService(sess)

	for _, v := range vlans {
		realvlan, err := networkVlanService.Mask(objectMask).Id(*v.Id).GetObject()
		if err != nil {
			log.Fatal(err)
		}

		for _, s := range realvlan.Subnets {

			if err != nil {
				log.Fatal(err)
			}

			if strings.Contains(*s.SubnetType, "PRIMARY") && !strings.Contains(*s.SubnetType, "_6") {

				templateVlan := Vlan{VlanId: *realvlan.VlanNumber}

				realSubnet, err := networkSubnetService.Mask(subnetObjectMask).Id(*s.Id).GetObject()

				if err != nil {
					log.Fatal(err)
				}

				ipAddresses := make([]string, 0, *realSubnet.IpAddressCount)

				for _, ip := range realSubnet.IpAddresses {
					ipAddresses = append(ipAddresses, *ip.IpAddress)
				}

				subnetvlanmap[*realvlan.VlanNumber] = Subnet{
					Cidr:               *realSubnet.Cidr,
					DnsServer:          *realSubnet.Gateway,
					MachineNetworkCidr: fmt.Sprintf("%s/%d", *realSubnet.NetworkIdentifier, *realSubnet.Cidr),
					Gateway:            *realSubnet.Gateway,
					Mask:               *realSubnet.Netmask,
					Network:            *realSubnet.NetworkIdentifier,
					IpAddresses:        ipAddresses,
				}
				templateVlan.Subnet = subnetvlanmap[*realvlan.VlanNumber]

				err = intTemplate.Execute(os.Stdout, templateVlan)
				if err != nil {
					log.Fatal(err)
				}
				err = pgTemplate.Execute(os.Stdout, templateVlan)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	jsonResults, err := json.MarshalIndent(subnetvlanmap, "", "    ")

	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile("subnets.json", jsonResults, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

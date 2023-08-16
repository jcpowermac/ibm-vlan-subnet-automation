package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	_ "strings"

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
	DhcpEndLocation    int
}

func main() {

	/*
	 * $notes = "gateway: %s, cidr: %d, vlan: %d"
	 * New-VDPortgroup -Name "ci-vlan-%d" -Notes $notes -VDSwitch $vdswitch -VLanId %d
	 */

	powercli := `
$vdswitchName = "vcs8e-vcs-ci-workload-private"
$vdswitch = Get-VDSwitch -Name $vdswitchName 
$notes = "vlan: {{.VlanId}} gateway: {{.Subnet.Gateway}} cidr: {{.Subnet.Cidr}} mask: {{.Subnet.Mask}}"
New-VDPortgroup -Name "ci-vlan-{{.VlanId}}" -Notes $notes -VDSwitch $vdswitch -VLanId {{.VlanId}}
`

	interfaces := `
interfaces {
	bonding dp0bond0 {
		vif {{.VlanId}} {
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
}
service {
	dhcp-server {
		shared-network-name ci-vlan-{{.VlanId}} {
			subnet {{.Subnet.Network}}/{{.Subnet.Cidr}} {
				default-router {{.Subnet.Gateway}}
				dns-server {{.Subnet.Gateway}}
				lease 3600
				ping-check
				start {{ index .Subnet.IpAddresses 10 }} {
					stop {{index .Subnet.IpAddresses (.Subnet.DhcpEndLocation) }}
				}
			}
		}
	}
	dns {
		forwarding {
			listen-on dp0bond0.{{.VlanId}} 
		}
	}
	nat {
		source {
			rule XXXX {
                outbound-interface dp0bond1
                source {
					address {{.Subnet.Network}}/{{.Subnet.Cidr}}
                }
                translation {
					address masquerade
				}
			}
		}
	}
}
`

	vlansFile, err := os.Open("configurations/vlans.txt")

	if err != nil {
		log.Fatal()
	}

	newVlans := make(map[int]struct{})
	scanner := bufio.NewScanner(vlansFile)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		tempInt, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		newVlans[tempInt] = struct{}{}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

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

			//if strings.Contains(*s.SubnetType, "PRIMARY") && !strings.Contains(*s.SubnetType, "_6") {
			if _, ok := newVlans[*realvlan.VlanNumber]; ok {

				vyattaConfFile, err := os.Create(fmt.Sprintf("configurations/vyatta-%d-01.conf", *realvlan.VlanNumber))

				if err != nil {
					log.Fatal(err)
				}
				powerShellFile, err := os.Create(fmt.Sprintf("configurations/vlan-%d.ps1", *realvlan.VlanNumber))
				if err != nil {
					log.Fatal(err)
				}
				//vyattaWriter := bufio.NewWriter(vyattaConfFile)

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
					DhcpEndLocation:    len(ipAddresses) - 2,
				}
				templateVlan.Subnet = subnetvlanmap[*realvlan.VlanNumber]

				err = intTemplate.Execute(vyattaConfFile, templateVlan)
				if err != nil {
					log.Fatal(err)
				}
				err = vyattaConfFile.Close()
				if err != nil {
					log.Fatal(err)
				}

				err = pgTemplate.Execute(powerShellFile, templateVlan)
				if err != nil {
					log.Fatal(err)
				}
				err = powerShellFile.Close()
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

	err = os.WriteFile("configurations/subnets.json", jsonResults, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

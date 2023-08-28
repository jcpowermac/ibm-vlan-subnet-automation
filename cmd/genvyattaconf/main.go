package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"os"
	"path/filepath"
	"strconv"
	_ "strings"

	"github.com/c-robinson/iplib"
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
	VirtualCenter      string   `json:"virtualcenter"`
	IPv6Prefix         string   `json:"ipv6prefix"`

	VifIpAddress string

	DhcpEndLocation int
}

func generateSlash30(vlanId uint32) (net.IP, net.IP) {
	var generatedVlanId uint32
	classC := net.ParseIP("192.168.0.0")
	ipInt := iplib.IP4ToUint32(classC)

	// Using just the vlanid int or *2
	// is not enough to stop collisions
	generatedVlanId = vlanId * 4
	newIpInt := ipInt + generatedVlanId
	generatedIpInt := iplib.Uint32ToIP4(newIpInt)

	newSlash30 := iplib.NewNet4(generatedIpInt, 30)

	first := newSlash30.FirstAddress()

	last := newSlash30.LastAddress()

	return first, last
}

func main() {

	powercli := `
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: {{.VlanId}} gateway: {{.Subnet.Gateway}} cidr: {{.Subnet.Cidr}} mask: {{.Subnet.Mask}}"
$newpg = New-VDPortgroup -Name "ci-vlan-{{.VlanId}}" -Notes $notes -VDSwitch $vdswitch -VLanId {{.VlanId}}

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True}
`

	vyattaConfTemplate := `
interfaces {
	bonding dp0bond0 {
		vif {{.VlanId}} {
			address {{.VifIpAddress}}/30
			vrrp {
				vrrp-group 1 {
					preempt false
					priority 254
					sync-group vgroup1
					virtual-address {{.Subnet.Gateway}}/{{.Subnet.Cidr}}
				}
			}
			ipv6 {
				router-advert {
					default-preference high
					managed-flag true
					max-interval 30 
					min-interval 30 
					prefix {{.Subnet.IPv6Prefix}} {
       					valid-lifetime 2592000
   					}
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
			rule {{.VlanId}} {
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

	virtualCenter := "vcs8e-vc.ocp2.dev.cluster.com"

	vlansAbsPath, err := filepath.Abs("../../configurations/vlans.txt")

	if err != nil {
		log.Fatal(err)
	}

	vlansFile, err := os.Open(vlansAbsPath)

	if err != nil {
		log.Fatal()
	}

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
	intTemplate, err := template.New("interface").Parse(vyattaConfTemplate)

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

	// fd65:a1a8:60ad:271c::1/64,fd65:a1a8:60ad:271c::1
	//

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
				ipv6Subnet := iplib.Net6FromStr(fmt.Sprintf("fd65:a1a8:60ad:%d::1/64", *realvlan.VlanNumber))

				vyattaConfFile01, err := os.Create(fmt.Sprintf("../../configurations/vyatta-%d-01.conf", *realvlan.VlanNumber))

				if err != nil {
					log.Fatal(err)
				}
				vyattaConfFile02, err := os.Create(fmt.Sprintf("../../configurations/vyatta-%d-02.conf", *realvlan.VlanNumber))

				if err != nil {
					log.Fatal(err)
				}
				powerShellFile, err := os.Create(fmt.Sprintf("../../configurations/vlan-%d.ps1", *realvlan.VlanNumber))
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

				first, last := generateSlash30(uint32(*realvlan.VlanNumber))

				subnetvlanmap[*realvlan.VlanNumber] = Subnet{
					Cidr:               *realSubnet.Cidr,
					DnsServer:          *realSubnet.Gateway,
					MachineNetworkCidr: fmt.Sprintf("%s/%d", *realSubnet.NetworkIdentifier, *realSubnet.Cidr),
					Gateway:            *realSubnet.Gateway,
					Mask:               *realSubnet.Netmask,
					Network:            *realSubnet.NetworkIdentifier,
					IpAddresses:        ipAddresses,
					DhcpEndLocation:    len(ipAddresses) - 2,
					VirtualCenter:      virtualCenter,

					IPv6Prefix: ipv6Subnet.String(),

					VifIpAddress: first.String(),
				}
				templateVlan.Subnet = subnetvlanmap[*realvlan.VlanNumber]

				err = intTemplate.Execute(vyattaConfFile01, templateVlan)
				if err != nil {
					log.Fatal(err)
				}
				err = vyattaConfFile01.Close()
				if err != nil {
					log.Fatal(err)
				}

				templateVlan.VifIpAddress = last.String()

				err = intTemplate.Execute(vyattaConfFile02, templateVlan)
				if err != nil {
					log.Fatal(err)
				}
				err = vyattaConfFile02.Close()
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

	err = os.WriteFile("../../configurations/subnets.json", jsonResults, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

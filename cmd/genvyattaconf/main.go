package main

import (
	_ "bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net"
	"os"
	_ "path/filepath"
	_ "strconv"
	"strings"
	_ "strings"

	"github.com/c-robinson/iplib"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

const (
	powercli = `
#$vdswitchName = "vcs8e-vcs-ci-workload-private"
# $pgrolename = "openshift_portgroup"
# $principal = "VSPHERE.LOCAL\openshift"

$pgrole = get-virole $pgrolename


$vdswitch = Get-VDSwitch -Name $vdswitchName
$notes = "vlan: {{.VlanId}} gateway: {{.Subnet.Gateway}} cidr: {{.Subnet.Cidr}} mask: {{.Subnet.Mask}}"
$newpg = New-VDPortgroup -Name "ci-vlan-{{.VlanId}}" -Notes $notes -VDSwitch $vdswitch -VLanId {{.VlanId}}

New-VIPermission -Entity $newpg -Principal $principal -Role $pgrole -Propagate $True
`

	vyattaConfTemplate = `
interfaces {
	bonding dp0bond0 {
		vif {{.VlanId}} {
			address {{.VifIpAddress}}/30
			address {{.VifIPv6Address}}/{{.Subnet.CidrIPv6}}
			vrrp {
				vrrp-group 1 {
					preempt false
					priority {{.Subnet.Priority}} 
					sync-group vgroup1
					virtual-address {{.Subnet.Gateway}}/{{.Subnet.Cidr}}
				}
				vrrp-group 2 {
					preempt false
					priority {{.Subnet.Priority}} 
					sync-group vgroup1
					virtual-address {{.Subnet.GatewayIPv6}}/{{.Subnet.CidrIPv6}}
					virtual-address {{.Subnet.LinkLocalIPv6}}
				}
			}
			ipv6 {
				router-advert {
					default-preference high
					managed-flag true
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
	dhcpv6-server {
		shared-network-name ci-vlan-{{.VlanId}} {
			subnet {{.Subnet.IPv6Prefix}} {
				address-range {
					start {{ .Subnet.StartIPv6Address }} {
						stop {{ .Subnet.StopIPv6Address }}
					}
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
)

type Vlan struct {
	VlanId int `json:"vlanId"`
	Subnet `json:"subnet"`
}

type Subnet struct {
	Cidr               int `json:"cidr"`
	CidrIPv6           int
	DnsServer          string   `json:"dnsServer"`
	MachineNetworkCidr string   `json:"machineNetworkCidr"`
	Gateway            string   `json:"gateway"`
	GatewayIPv6        string   `json:"gatewayipv6"`
	Mask               string   `json:"mask"`
	Network            string   `json:"network"`
	IpAddresses        []string `json:"ipAddresses"`
	VirtualCenter      string   `json:"virtualcenter"`
	IPv6Prefix         string   `json:"ipv6prefix"`
	StartIPv6Address   string
	StopIPv6Address    string
	LinkLocalIPv6      string

	VifIpAddress   string
	VifIPv6Address string

	DhcpEndLocation int

	Priority int
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

	virtualCenter := "vcenter.ibmc.devcluster.openshift.com"
	//primaryRouterHostname := "bcr01a.dal10"
	//subnetvlanmap := make(map[int]Subnet)

	subnetvlanmap := make(map[string]map[int]Subnet)

	/*
		vlansAbsPath, err := filepath.Abs("../../configurations/vlans.txt")

		if err != nil {
			log.Fatal(err)
		}

		vlansFile, err := os.Open(vlansAbsPath)

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

	*/

	pgTemplate, err := template.New("pg").Parse(powercli)

	if err != nil {
		log.Fatal(err)
	}
	intTemplate, err := template.New("interface").Parse(vyattaConfTemplate)

	if err != nil {
		log.Fatal(err)
	}

	// https://github.com/softlayer/softlayer-go/blob/master/examples/cmd/vlan_detail.go
	// https://sldn.softlayer.com/reference/datatypes/SoftLayer_Network_Subnet/

	objectMask :=
		`mask[id,name,vlanNumber,subnets[id,cidr,netmask,networkIdentifier,subnetType],primaryRouter[hostname]]`
	subnetObjectMask :=
		`mask[id,ipAddressCount,gateway,cidr,netmask,networkIdentifier,ipAddresses[ipAddress,isNetwork,isBroadcast,isGateway]]`

	sess := session.New()

	service := services.GetAccountService(sess)
	networkVlanService := services.GetNetworkVlanService(sess)

	vlans, err := service.Mask("mask[id,name,vlanNumber,primaryRouter[hostname]]").GetNetworkVlans()

	if err != nil {
		log.Fatal(err)
	}

	networkSubnetService := services.GetNetworkSubnetService(sess)

	for _, v := range vlans {
		if v.Name != nil {
			if strings.Contains(*v.Name, "ci") {
				realvlan, err := networkVlanService.Mask(objectMask).Id(*v.Id).GetObject()
				if err != nil {
					log.Fatal(err)
				}

				for _, s := range realvlan.Subnets {
					if !strings.Contains(*s.SubnetType, "GLOBAL") {
						ipv6Subnet := iplib.Net6FromStr(fmt.Sprintf("fd65:a1a8:60ad:%d::1/64", *realvlan.VlanNumber))
						linkLocalSubnet := iplib.Net6FromStr("fe80::/64")

						cidrLinkLocalIPv6, _ := linkLocalSubnet.Mask().Size()

						cidripv6, _ := ipv6Subnet.Mask().Size()

						err := os.MkdirAll(fmt.Sprintf("../../configurations/%s", *realvlan.PrimaryRouter.Hostname), 0750)

						if err != nil {
							log.Fatal(err)
						}
						vyattaConfFile01, err := os.Create(fmt.Sprintf("../../configurations/%s/vyatta-%d-01.conf", *realvlan.PrimaryRouter.Hostname, *realvlan.VlanNumber))

						if err != nil {
							log.Fatal(err)
						}
						vyattaConfFile02, err := os.Create(fmt.Sprintf("../../configurations/%s/vyatta-%d-02.conf", *realvlan.PrimaryRouter.Hostname, *realvlan.VlanNumber))

						if err != nil {
							log.Fatal(err)
						}
						powerShellFile, err := os.Create(fmt.Sprintf("../../configurations/%s/vlan-%d.ps1", *realvlan.PrimaryRouter.Hostname, *realvlan.VlanNumber))
						if err != nil {
							log.Fatal(err)
						}

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

						if _, ok := subnetvlanmap[*realvlan.PrimaryRouter.Hostname]; !ok {
							subnetvlanmap[*realvlan.PrimaryRouter.Hostname] = make(map[int]Subnet)
						}

						subnetvlanmap[*realvlan.PrimaryRouter.Hostname][*realvlan.VlanNumber] = Subnet{
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

							GatewayIPv6: ipv6Subnet.Enumerate(1, 2)[0].String(),

							StartIPv6Address: ipv6Subnet.Enumerate(1, 4)[0].String(),
							StopIPv6Address:  ipv6Subnet.Enumerate(1, 100)[0].String(),
							LinkLocalIPv6:    fmt.Sprintf("%s/%d", linkLocalSubnet.Enumerate(1, *realvlan.VlanNumber)[0].String(), cidrLinkLocalIPv6),

							CidrIPv6: cidripv6,

							VifIPv6Address: ipv6Subnet.Enumerate(1, 1)[0].String(),

							VifIpAddress: first.String(),
							Priority:     254,
						}
						templateVlan.Subnet = subnetvlanmap[*realvlan.PrimaryRouter.Hostname][*realvlan.VlanNumber]

						err = intTemplate.Execute(vyattaConfFile01, templateVlan)
						if err != nil {
							log.Fatal(err)
						}
						err = vyattaConfFile01.Close()
						if err != nil {
							log.Fatal(err)
						}

						templateVlan.Subnet.VifIpAddress = last.String()
						templateVlan.Subnet.VifIPv6Address = ipv6Subnet.Enumerate(1, 3)[0].String()
						templateVlan.Subnet.Priority = 253

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

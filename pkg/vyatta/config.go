package vyatta

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

const (
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
)



func GenerateConfig() {

	intTemplate, err := template.New("interface").Parse(vyattaConfTemplate)

	if err != nil {
		log.Fatal(err)
	}
	vyattaConfFile01, err := os.Create(fmt.Sprintf("../../configurations/vyatta-%d-01.conf", *realvlan.VlanNumber))

	if err != nil {
		log.Fatal(err)
	}
	vyattaConfFile02, err := os.Create(fmt.Sprintf("../../configurations/vyatta-%d-02.conf", *realvlan.VlanNumber))

	if err != nil {
		log.Fatal(err)
	}

}
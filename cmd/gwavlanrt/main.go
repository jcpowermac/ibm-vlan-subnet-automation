package main

import (
	"bufio"
	"github.com/softlayer/softlayer-go/sl"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

const (
	hardwareMask      = `mask[id,hostname,domain,networkComponents[uplinkComponent[id,networkVlan,networkVlanTrunks[networkVlan[id,vlanNumber]]]]]`
	gatewayMemberMask = `mask[id,hardware[networkComponents[uplinkComponent[id,networkVlan,networkVlanTrunks[networkVlan[id,vlanNumber]]]]]]`
)

func main() {

	gatewayId, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	vlanMap := make(map[int]datatypes.Network_Vlan)
	var newVlansToTrunk []datatypes.Network_Vlan

	var newVlansToUnbypass []datatypes.Network_Gateway_Vlan
	sess := session.New()
	service := services.GetAccountService(sess)

	//networkGateway := services.GetNetworkGatewayService(sess)

	/*
		hardware, err := service.Mask(hardwareMask).GetHardware()
		if err != nil {
			log.Fatal(err)
		}

	*/
	vlans, err := service.Mask("mask[id,name,vlanNumber,primaryRouter[hostname]]").GetNetworkVlans()
	if err != nil {
		log.Fatal(err)
	}

	primaryRouterHostname := "bcr01a.dal10"
	//networkGatewayService := services.GetNetworkGatewayService(sess)

	for _, v := range vlans {
		if *v.PrimaryRouter.Hostname == primaryRouterHostname {
			vlanMap[*v.VlanNumber] = v
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	vlansAbsPath, err := filepath.Abs("../../configurations/vlans.txt")

	if err != nil {
		log.Fatal(err)
	}

	vlansFile, err := os.Open(vlansAbsPath)

	if err != nil {
		log.Fatal()
	}
	scanner := bufio.NewScanner(vlansFile)

	for scanner.Scan() {
		tempInt, err := strconv.Atoi(scanner.Text())

		tempNetworkVlan := vlanMap[tempInt]

		if err != nil {
			log.Fatal(err)
		}

		newVlansToTrunk = append(newVlansToTrunk, tempNetworkVlan)

		newVlansToUnbypass = append(newVlansToUnbypass, datatypes.Network_Gateway_Vlan{NetworkVlan: &tempNetworkVlan, NetworkVlanId: tempNetworkVlan.Id})

		log.Printf("vlan %d, id %d", &tempNetworkVlan.VlanNumber, &tempNetworkVlan.Id)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	networkGateways, err := service.GetNetworkGateways()
	if err != nil {
		log.Fatal(err)
	}
	networkGatewayVlanService := services.GetNetworkGatewayVlanService(sess)

	for _, ng := range networkGateways {
		log.Printf("%s %d", *ng.Name, *ng.Id)

		if gatewayId == *ng.Id {

			for _, v := range newVlansToTrunk {
				templateObject := datatypes.Network_Gateway_Vlan{
					BypassFlag:       sl.Bool(false),
					NetworkGatewayId: ng.Id,
					NetworkVlanId:    v.Id,
				}

				_, err = networkGatewayVlanService.CreateObject(&templateObject)
				if err != nil {
					log.Print(err)
				}
			}

		}
	}
}

/*
	members, err := networkGateway.Mask(gatewayMemberMask).GetMembers()

	if err != nil {
		log.Fatal(err)
	}

*/

/*
	insideVlans, err := networkGateway.Mask("mask[bypassFlag,id,networkVlanId,networkVlan[id,vlanNumber]]").GetInsideVlans()

	if err != nil {
		log.Fatal(err)
	}

	for _, iv := range insideVlans {
		log.Printf("vlan id: %d, gateway vlan id: %d, vlan number %d", *iv.NetworkVlanId, *iv.Id, *iv.NetworkVlan.VlanNumber)
	}


	for _, m := range members {

		log.Printf("member id: %d", *m.Id)
		log.Printf("hardware id: %d", *m.Hardware.Id)

		primaryBackendNetworkComponent, err := services.GetHardwareService(sess).Id(*m.Hardware.Id).GetPrimaryBackendNetworkComponent()

		if err != nil {
			log.Fatal(err)
		}
		log.Printf("primary id %d, name %s", *primaryBackendNetworkComponent.Id, *primaryBackendNetworkComponent.Name)
		networkComponent := services.GetNetworkComponentService(sess).Id(*primaryBackendNetworkComponent.Id)
		_, err = networkComponent.AddNetworkVlanTrunks(newVlansToTrunk)

		for _, nc := range m.Hardware.NetworkComponents {
			for _, t := range nc.UplinkComponent.NetworkVlanTrunks {
				log.Printf("trunk vlan id %d", *t.NetworkVlan.VlanNumber)
			}
		}

	}

	err = networkGateway.UnbypassVlans(newVlansToUnbypass)
	if err != nil {
		log.Fatal(err)
	}
*/

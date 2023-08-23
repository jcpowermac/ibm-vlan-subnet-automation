package main

import (
	"bufio"
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

	vlanMap := make(map[int]datatypes.Network_Vlan)
	var newVlansToTrunk []datatypes.Network_Vlan
	sess := session.New()
	service := services.GetAccountService(sess)

	//networkGateway := services.GetNetworkGatewayService(sess)

	/*
		hardware, err := service.Mask(hardwareMask).GetHardware()
		if err != nil {
			log.Fatal(err)
		}

	*/
	vlans, err := service.GetNetworkVlans()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range vlans {
		vlanMap[*v.VlanNumber] = v
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

		newVlanToTrunk := vlanMap[tempInt]

		if err != nil {
			log.Fatal(err)
		}
		newVlansToTrunk = append(newVlansToTrunk, datatypes.Network_Vlan{VlanNumber: newVlanToTrunk.VlanNumber, Id: newVlanToTrunk.Id})
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	networkGateways, err := service.GetNetworkGateways()
	if err != nil {
		log.Fatal(err)
	}

	for _, ng := range networkGateways {
		log.Printf("%s %d", *ng.Name, *ng.Id)

		networkGatewayService := services.GetNetworkGatewayService(sess).Id(*ng.Id)

		members, err := networkGatewayService.Mask(gatewayMemberMask).GetMembers()

		if err != nil {
			log.Fatal(err)
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

		networkGatewayService.UnbypassVlans()
	}

	/*
		for _, h := range hardware {

			//log.Printf("%s", *h.Hostname)

			for _, nc := range h.NetworkComponents {
				for _, t := range nc.UplinkComponent.NetworkVlanTrunks {
					log.Printf("trunk vlan id %d", *t.NetworkVlan.VlanNumber)
				}
			}
		}

	*/
}

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

func main() {
	sess := session.New()
	service := services.GetAccountService(sess)

	//networkService := services.GetNetworkService(sess)
	hardwareMask := `mask[id,hostname,domain,networkComponents[uplinkComponent[id,networkVlan,networkVlanTrunks[networkVlan[id,vlanNumber]]]]]`
	hardware, err := service.Mask(hardwareMask).GetHardware()
	if err != nil {
		log.Fatal(err)
	}

	vlanMap := make(map[int]datatypes.Network_Vlan)

	vlans, err := service.GetNetworkVlans()
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range vlans {
		vlanMap[*v.VlanNumber] = v
	}

	//networkComponent := services.GetNetworkComponentService(sess)

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

	var newVlansToTrunk []datatypes.Network_Vlan
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

	for _, h := range hardware {
		log.Printf("name %s id %d", *h.Hostname, *h.Id)
		if *h.Id == 3031402 {
			primaryBackendNetworkComponent, err := services.GetHardwareService(sess).Id(*h.Id).GetPrimaryBackendNetworkComponent()

			if err != nil {
				log.Fatal(err)
			}
			log.Printf("primary id %d, name %s", *primaryBackendNetworkComponent.Id, *primaryBackendNetworkComponent.Name)

			networkComponent := services.GetNetworkComponentService(sess).Id(*primaryBackendNetworkComponent.Id)

			log.Printf("testing ... %v", networkComponent)

			_, err = networkComponent.AddNetworkVlanTrunks(newVlansToTrunk)

			if err != nil {
				log.Fatal(err)
			}
		}

		/*
			for _, nc := range h.NetworkComponents {
				log.Printf("uplinke component id %d", nc.UplinkComponent.Id)

				//netComp := networkComponent.Id(*nc.Id)

				log.Printf("%s, %d", *nc.Name, *nc.Id)

				for _, cv := range nc.UplinkComponent.NetworkVlanTrunks {
					log.Printf("trunk vlan: %d", *cv.NetworkVlan.VlanNumber)
				}
				//netComp.Options.Mask = `mask[networkVlan[id,vlanNumber]]`

					currentTrunks, err := netComp.GetNetworkVlanTrunks()
					if err != nil {
						log.Fatal(err)
					}

					for _, vlans := range currentTrunks {
						log.Printf("vlan %d", *vlans.NetworkVlanId)
					}

			}
		*/
	}
}

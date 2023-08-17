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
	hardwareMask = `mask[id,hostname,domain,networkComponents[uplinkComponent[id,networkVlan,networkVlanTrunks[networkVlan[id,vlanNumber]]]]]`
)

func main() {
	hardwareIds := make(map[int]struct{})
	vlanMap := make(map[int]datatypes.Network_Vlan)
	var newVlansToTrunk []datatypes.Network_Vlan

	sess := session.New()
	service := services.GetAccountService(sess)

	hardware, err := service.Mask(hardwareMask).GetHardware()
	if err != nil {
		log.Fatal(err)
	}

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

	hardwareAbsPath, err := filepath.Abs("../../configurations/hardware.txt")
	if err != nil {
		log.Fatal(err)
	}

	hardwareFile, err := os.Open(hardwareAbsPath)

	if err != nil {
		log.Fatal()
	}
	hwscan := bufio.NewScanner(hardwareFile)

	for hwscan.Scan() {
		tempInt, err := strconv.Atoi(hwscan.Text())
		if err != nil {
			log.Fatal(err)
		}
		hardwareIds[tempInt] = struct{}{}

	}
	if err := hwscan.Err(); err != nil {
		log.Fatal(err)
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

	for _, h := range hardware {
		log.Printf("name %s id %d", *h.Hostname, *h.Id)
		if _, ok := hardwareIds[*h.Id]; ok {
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
	}
}

package main

import (
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
)

func main() {

	sess := session.New()

	service := services.GetAccountService(sess)

	//networkService := services.GetNetworkService(sess)
	service.GetHardware()

}

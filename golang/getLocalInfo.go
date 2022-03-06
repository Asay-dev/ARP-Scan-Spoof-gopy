package main

import (
	"net"
	"os"
	"strings"

	"github.com/withmandala/go-log"
)

func getLocalInfo(interfaceName string) (string, string) {

	logger := log.New(os.Stderr).WithColor()
	//----------------------
	// Get the local machine IP address
	// https://www.socketloop.com/tutorials/golang-how-do-I-get-the-local-ip-non-loopback-address
	//----------------------

	//addrs, err := net.InterfaceAddrs()

	//if err != nil {
	//	fmt.Println(err)
	//}

	//var currentIP string
	var loclaMAC string
	var localIP string

	//for _, address := range addrs {
	//
	//	// check the address type and if it is not a loopback the display it
	//	// = GET LOCAL IP ADDRESS
	//	if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	//		if ipnet.IP.To4() != nil {
	//			fmt.Println("Current IP address : ", ipnet.IP.String())
	//			currentIP = ipnet.IP.String()
	//		}
	//	}
	//}

	logger.Info("------------------------------")
	logger.Info("We want the interface name that has the current IP address")
	logger.Info("MUST NOT be binded to 127.0.0.1 ")
	logger.Info("------------------------------")

	// get all the system's or local machine's network interfaces

	interfaces, _ := net.Interfaces()
	for _, interf := range interfaces {

		if addrs, err := interf.Addrs(); err == nil {
			for index, addr := range addrs {
				logger.Info("[", index, "]", interf.Name, ">", addr)
				if interf.Name == interfaceName {
					localIP = strings.Split(addr.String(), "/")[0]
				}
			}
		}
	}

	logger.Info("------------------------------")

	// extract the hardware information base on the interface name
	// capture above
	//var localIP net.IP
	netInterface, err := net.InterfaceByName(interfaceName)

	if err != nil {
		logger.Error(err)
	}

	name := netInterface.Name
	loclaMAC = netInterface.HardwareAddr.String()

	logger.Warn("Hardware name : ", name)
	logger.Warn("IP address : ", localIP)

	logger.Warn("MAC address : ", loclaMAC)

	// verify if the MAC address can be parsed properly
	hwAddr, err := net.ParseMAC(loclaMAC)

	if err != nil {
		logger.Error("No able to parse MAC address : ", err)
		os.Exit(-1)
	}

	logger.Warn("Physical hardware address : %s \n", hwAddr.String())
	return localIP, loclaMAC

}

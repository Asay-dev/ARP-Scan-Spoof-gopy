package main

import (
	"ARTScript_ARP/arpScan"
	"log"
	"net"

	"github.com/jackpal/gateway"
)

func main() {
	// device := ""
	// gt := (net.ParseIP("192.168.31.1"))[12:]
	// tg := (net.ParseIP("192.168.31.96"))[12:]
	// lc := (net.ParseIP("192.168.31.254"))[12:]
	// tgm, _ := net.ParseMAC("")
	// gtm, _ := net.ParseMAC("")
	// lcm, _ := net.ParseMAC("")
	// arpPosion.ArpPoison2(device, gtm, gt, lcm, lc, tgm, tg)
	// arpScan.Start_scan()
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, iface := range ifaces {
		var addr *net.IPNet
		if addrs, err := iface.Addrs(); err != nil {
			return
		} else {
			for _, a := range addrs {
				if ipnet, ok := a.(*net.IPNet); ok {
					if ip4 := ipnet.IP.To4(); ip4 != nil {
						addr = &net.IPNet{
							// IP:   ip4,
							IP:   ip4,
							Mask: ipnet.Mask[len(ipnet.Mask)-4:],
						}
						break
					}
				}
			}
		}
		log.Printf("[-] interface %v:%v", iface.Name, addr)
	}

	if ip, err := gateway.DiscoverGateway(); err != nil {
		log.Printf("[!] ERROR : %v", err)
	} else {
		log.Printf("[*] Gateway : %v", ip)
		arpScan.Start_scan("WLAN", ip)
	}
	log.Println("==================================================================================================")

	arpScan.Start_scan2()

	// if ip, err := gateway.DiscoverInterface(); err != nil {
	// 	log.Printf("ERROR : %v", err)
	// } else {
	// 	log.Printf("Interface : %v", ip)
	// }
}

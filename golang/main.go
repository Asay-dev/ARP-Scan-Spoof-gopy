package main

import (
	"ARTScript_ARP/arpScan"
	"fmt"
	"log"
	"net"

	"github.com/google/gopacket/pcap"
	"github.com/jackpal/gateway"
)

var lc_ip net.IP
var interface_ip net.IP

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
	// 查找路由器
	if ip, err := gateway.DiscoverGateway(); err != nil {
		log.Printf("[!] ERROR : %v", err)
	} else {
		log.Printf("[*] Gateway : %v", ip)
		// arpScan.Start_scan("WLAN", ip)
	}
	// 查看网卡
	if lc_ip, err := gateway.DiscoverInterface(); err != nil {
		log.Printf("ERROR : %v", err)
	} else {
		log.Printf("Interface : %v", lc_ip)
	}

	log.Println("==================================================================================================")

	// 得到所有的(网络)设备
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	// 打印设备信息
	fmt.Println("Devices found:")
	for _, device := range devices {

		for _, address := range device.Addresses {
			if fmt.Sprintf("%s", address.IP) == "192.168.0.194" {
				fmt.Println("\nName: ", device.Name)
				fmt.Println("Description: ", device.Description)
				fmt.Println("Devices addresses: ", device.Description)

				fmt.Println("- IP address: ", address.IP)
				fmt.Println("- Subnet mask: ", address.Netmask)
				arpScan.Start_scan3(device.Name)
			}
		}
		// arpScan.Start_scan3(device.Name)

	}

}

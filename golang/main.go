package main

import (
	"encoding/hex"
	"os"
	"strings"
	"time"

	"github.com/HayatoDoi/arp-scan-X/arp"
	"github.com/HayatoDoi/arp-scan-X/syslog"
)

//go:generate oui -p main -o oui.go

func main() {

	Set_opts()
	slog := syslog.New(opts.DebugMode)

	interfaces, err := arp.IfaceToName("WLAN")
	if err != nil {
		slog.Errorln("%s", err)

		os.Exit(1)
	}

	success := false
	for _, interface_ := range interfaces {
		// make config
		config := arp.Config{
			Interface: interface_,
			//Timeout:   time.Duration(500) * time.Millisecond,
			Timeout: time.Duration(opts.Timeout) * time.Millisecond,
			//Backoff: 1.5,
			Backoff: opts.Backoff,
		}
		a, err := arp.New(config)
		if err != nil {
			slog.Debugln("Error(%s) : %s", interface_, err)
			continue
		}
		slog.Println("Interface: %s, Network range: %v", interface_, a.Addr)
		arpTables, err := a.Scan()
		if err != nil {
			slog.Debugln("Error(%s) : %s", interface_, err)
			continue
		}
		for _, arpTable := range arpTables {
			oui := strings.ToUpper(hex.EncodeToString(arpTable.HardwareAddr[:3]))
			organization, ok := MacAndOrganization[oui]
			if ok != true {
				organization = "unknown"
			}
			slog.Println("%-15v %-20v %s", arpTable.IP, arpTable.HardwareAddr, organization)
		}
		success = true
	}
	if success == false {
		slog.Errorln("No valid ip address configuration.")
		os.Exit(1)
	}
}

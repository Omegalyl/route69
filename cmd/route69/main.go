package main

import (
	"fmt"
	"route69/internal/capture"
	"route69/internal/parse"
	"strings"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("Starting Capture!!")
	capCh := make(chan *capture.Capture, 100)
	go capture.CapturePackets(capCh, "br-c31d9f6f8c43")
	for cap := range capCh {
		frame, err := parse.ParseEthernet2Frame(cap.Pk, cap.Addr.(*unix.SockaddrLinklayer))
		if err != nil {
			continue
		}

		fmt.Println(strings.Repeat("=", 50))
		switch frame.EtherType {
		case unix.ETH_P_ARP:
			arp, _ := parse.ParseARP(frame.Payload)
			fmt.Println("Ethernet:\n", frame)
			fmt.Println("ARP:\n", arp)
		default:
			fmt.Println("Ethernet:\n", frame)
		}
		fmt.Println(strings.Repeat("=", 50))
	}
}

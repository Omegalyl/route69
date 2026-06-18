package main

import (
	"fmt"
	"route69/internal/capture"
	"route69/internal/parse"

	"golang.org/x/sys/unix"
)

func main() {
	fmt.Println("Starting Firewall!!")
	capCh := make(chan *capture.Capture, 100)
	go capture.CapturePackets(capCh)
	for cap := range capCh {
		frame, err := parse.ParseEthernet2Frame(cap.Pk, cap.Addr.(*unix.SockaddrLinklayer))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(frame)
	}
}

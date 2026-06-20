package capture

import (
	"log"
	"net"
	"route69/internal/utils"

	"golang.org/x/sys/unix"
)

type Capture struct {
	Pk   []byte
	Addr unix.Sockaddr
}

func CapturePackets(capCh chan *Capture, ifName string) {
	fd, err := unix.Socket(
		unix.AF_PACKET,
		unix.SOCK_RAW,
		int(utils.Htons(unix.ETH_P_ALL)),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer unix.Close(fd)
	iface, err := net.InterfaceByName(ifName)
	if err != nil {
		log.Fatal("Interface not found:", err)
	}
	if err := unix.Bind(fd, &unix.SockaddrLinklayer{
		Protocol: utils.Htons(unix.ETH_P_ALL),
		Ifindex:  iface.Index,
	}); err != nil {
		log.Fatal(err)
	}
	buf := make([]byte, 65535)

	for {
		n, addr, err := unix.Recvfrom(fd, buf, 0)
		if err != nil {
			continue
		}
		pk := make([]byte, n)
		copy(pk, buf[:n])
		capCh <- &Capture{
			Pk:   pk,
			Addr: addr,
		}

	}
}

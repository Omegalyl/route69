package parse

import (
	"encoding/binary"
	"fmt"
	"net"

	"golang.org/x/sys/unix"
)

type Ethernet2Frame struct {
	IfIndex   int
	HaType    uint16
	PkType    uint8
	DstMac    net.HardwareAddr
	SrcMac    net.HardwareAddr
	EtherType uint16
	Payload   []byte
}

func (e *Ethernet2Frame) String() string {
	return fmt.Sprintf(
		"Interface: %d\n Dst Mac: %s\n Scr Mac: %s\n EtherType: %d\n Payload Size: %d\n",
		e.IfIndex, e.DstMac, e.SrcMac, e.EtherType, len(e.Payload),
	)
}

func ParseEthernet2Frame(pk []byte, sll *unix.SockaddrLinklayer) (*Ethernet2Frame, error) {
	if sll.Hatype != unix.ARPHRD_ETHER {
		return nil, fmt.Errorf("ethernet2: only ethernet2 frame supported")
	}
	if len(pk) < 14 {
		return nil, fmt.Errorf("ethernet2: invalid frame size")
	}
	frame := &Ethernet2Frame{
		IfIndex: sll.Ifindex,
		PkType:  sll.Pkttype,
		HaType:  sll.Hatype,
		DstMac:  net.HardwareAddr(pk[:6]),
		SrcMac:  net.HardwareAddr(pk[6:12]),
	}
	etherType := binary.BigEndian.Uint16(pk[12:14])
	switch etherType {
	case unix.ETH_P_IP, unix.ETH_P_IPV6, unix.ETH_P_ARP:
		// supported Ethernet II types
	default:
		return nil, fmt.Errorf("ethernet2: unsupported ethertype 0x%04x", etherType)
	}
	frame.EtherType = etherType
	if len(pk[14:]) > 1500 {
		return nil, fmt.Errorf("ethernet2: Jumbo+ frames not supported")
	}
	frame.Payload = append([]byte(nil), pk[14:]...)
	return frame, nil
}

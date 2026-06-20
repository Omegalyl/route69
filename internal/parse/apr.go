package parse

import (
	"encoding/binary"
	"fmt"
	"net"
)

type ARP struct {
	Hatype uint16
	Ptype  uint16
	HLen   uint8
	PLen   uint8
	Oper   uint16
	SHA    []byte
	SPA    []byte
	THA    []byte
	TPA    []byte
}

func (a *ARP) String() string {
	op := "unknown"
	switch a.Oper {
	case 1:
		op = "request"
	case 2:
		op = "reply"
	}
	return fmt.Sprintf(
		"Oper: %s\n Sender: %s (%s)\n Target: %s (%s)\n",
		op,
		net.IP(a.SPA), net.HardwareAddr(a.SHA),
		net.IP(a.TPA), net.HardwareAddr(a.THA),
	)
}

func ParseARP(pk []byte) (*ARP, error) {
	if len(pk) < 8 {
		return nil, fmt.Errorf("ARP: invalid payload")
	}
	arp := &ARP{
		Hatype: binary.BigEndian.Uint16((pk[:2])),
		Ptype:  binary.BigEndian.Uint16((pk[2:4])),
		HLen:   pk[4],
		PLen:   pk[5],
		Oper:   binary.BigEndian.Uint16(pk[6:8]),
	}

	hlen, plen := int(arp.HLen), int(arp.PLen)
	if len(pk[8:]) < 2*hlen+2*plen {
		return nil, fmt.Errorf("ARP: invalid payload")
	}

	offset := 8
	arp.SHA = append([]byte(nil), pk[offset:offset+hlen]...)
	offset += hlen
	arp.SPA = append([]byte(nil), pk[offset:offset+plen]...)
	offset += plen
	arp.THA = append([]byte(nil), pk[offset:offset+hlen]...)
	offset += hlen
	arp.TPA = append([]byte(nil), pk[offset:offset+plen]...)

	return arp, nil
}

package ip

import (
	"net"
)

// IPRangeOfCIDR returns the start IP and end IP represents by the cidr
func IPV4RangeOfCIDR(cidr string) ([]net.IP, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	//calc wildcard mask
	ones, bits := ipnet.Mask.Size()

	ones = bits - ones
	wcMask := make(net.IPMask, 4)
	for i := 3; i >= 0; i-- {
		if ones >= 8 {
			wcMask[i] = 0xff
			ones -= 8
			continue
		}

		wcMask[i] = byte(0xff >> (8 - uint(ones)))
		ones = 0
	}

	endIP := make(net.IP, 4)
	for i := 0; i < 4; i++ {
		endIP[i] = ipnet.IP[i] | wcMask[i]
	}

	return []net.IP{ipnet.IP, endIP}, nil
}

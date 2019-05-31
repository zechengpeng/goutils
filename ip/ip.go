package ip

import (
	"net"
)

var privateIPBlocks []*net.IPNet

func init() {
	for _, cidr := range []string{
		"127.0.0.0/8",    // IPv4 loopback
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"::1/128",        // IPv6 loopback
		"fe80::/10",      // IPv6 link-local
		"fc00::/7",       // IPv6 unique local addr
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}
}

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

func IsPrivateIP(ip net.IP) bool {
	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}

	return false
}

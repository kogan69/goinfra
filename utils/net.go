package utils

import (
	"math/rand"
	"net"
)

// AddressToIp check if address is a valid ip. If yes, it returns it
// Otherwise, it will handle the address as a host name and will perform DNS lookup
func AddressToIp(address string) (ip string, err error) {
	addr := net.ParseIP(address)
	if addr != nil {
		ip = addr.String()
		return
	} else {
		// try as host
		var ips []string
		ips, err = net.LookupHost(address)
		if err != nil {
			return "", err
		} else {
			ip = ips[0]
			return
		}
	}
}

func IpToUint(ip net.IP) uint {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	return v
}

func UintToIpV4(n uint) net.IP {
	var b [4]byte
	b[0] = byte(n & 0xFF)
	b[1] = byte((n >> 8) & 0xFF)
	b[2] = byte((n >> 16) & 0xFF)
	b[3] = byte((n >> 24) & 0xFF)
	return net.IPv4(b[3], b[2], b[1], b[0])
}

func NumberOfHosts(cidr *net.IPNet) uint {
	ones, _ := cidr.Mask.Size()
	return 2 << (32 - ones - 1)
}
func Increment(cidr *net.IPNet, inc uint) *net.IP {
	if cidr == nil {
		return nil
	}
	v := IpToUint(cidr.IP.Mask(cidr.Mask))
	v += inc
	ip := UintToIpV4(v)
	if !cidr.Contains(ip) {
		return nil
	}
	return &ip
}

func RandomAddress(cidr *net.IPNet) *net.IP {
	hosts := int32(NumberOfHosts(cidr))
	inc := rand.Int31n(hosts - 1)
	return Increment(cidr, uint(inc))
}

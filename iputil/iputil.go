package iputil

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
)

func getIpRange(ipv4Net *net.IPNet) (uint32, uint32) {
	// convert IPNet struct mask and address to uint32
	// network is BigEndian
	mask := binary.BigEndian.Uint32(ipv4Net.Mask)
	start := binary.BigEndian.Uint32(ipv4Net.IP)

	// find the final address
	finish := (start & mask) | (mask ^ 0xffffffff)

	return start, finish
}
func getRandomIpUint32Array(start, finish uint32) []uint32 {
	// loop through addresses as uint32
	addresses := make([]uint32, 0, finish-start+1)
	for i := start; i <= finish; i++ {
		addresses = append(addresses, i)
	}

	// Shuffle the addresses
	rand.Shuffle(len(addresses), func(i, j int) { addresses[i], addresses[j] = addresses[j], addresses[i] })
	return addresses
}
func getValidIpArray(uint32arr []uint32, myip net.IP) []net.IP {
	ipArr := make([]net.IP, 0, len(uint32arr))
	for _, v := range uint32arr {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, v)
		if !ip.Equal(myip) && ip[3] != 0 && ip[3] != 255 {
			ipArr = append(ipArr, ip)
		}
	}
	return ipArr
}

func GetScanIpArray(cidr string) []net.IP {
	myip, ipv4Net, err := net.ParseCIDR(cidr)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	start, finish := getIpRange(ipv4Net)

	uint32arr := getRandomIpUint32Array(start, finish)
	ipArr := getValidIpArray(uint32arr, myip)
	return ipArr
}

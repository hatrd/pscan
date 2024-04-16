package ui

import (
	"fmt"
	"net"
)

func GetSubNet() net.Addr {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	for i, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			fmt.Println(i, ":", ipnet.String())
		}
	}
	fmt.Print("Select the subnet: ")
	var sel int
	_, err = fmt.Scan(&sel)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	fmt.Println(addrs[sel])
	return addrs[sel]
}

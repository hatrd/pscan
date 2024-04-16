package scan

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/hatrd/pscan/iputil"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func isAlive(target net.IP) bool {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Println("Error listening for ICMP packets: ", err)
		return false
	}
	defer conn.Close()

	message := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte(""),
		},
	}

	b, err := message.Marshal(nil)
	if err != nil {
		fmt.Println("Error marshaling ICMP message: ", target, err)
		return false
	}

	_, err = conn.WriteTo(b, &net.IPAddr{IP: target})
	if err != nil {
		fmt.Println("Error writing to connection: ", target, err)
		return false
	}

	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(icmp_dealine))
	if err != nil {
		fmt.Println("Error setting read deadline: ", target, err)
		return false
	}

	n, _, err := conn.ReadFrom(reply)
	if err != nil {
		fmt.Println("Error reading from connection: ", target, err)
		return false
	}

	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), reply[:n])
	if err != nil {
		fmt.Println("Error parsing ICMP message: ", err)
		return false
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Println("Host is alive: ", target.String())
		return true
	default:
		// fmt.Println("Host is not alive: ", target.String())
		return false
	}
}
func isPortOpen(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, port_deadline)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
func scanPort(ip string, port int, ch chan string) {
	if isPortOpen(ip, port) {
		ch <- fmt.Sprintf("%s:%d", ip, port)
		file, err := os.OpenFile(output_file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		_, err = file.WriteString(fmt.Sprintf("%s:%d\n", ip, port))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
func scanIp(ip net.IP, ch chan string) {

	cnt := make(chan int, port_concurrent)
	var wg sync.WaitGroup
	for _, port := range ports {
		cnt <- 1
		wg.Add(1)
		go func() {
			scanPort(ip.String(), port, ch)
			<-cnt
			wg.Done()
		}()
	}
	wg.Wait()
}
func Scan(addr net.Addr, ch chan string) {
	addresses := iputil.GetScanIpArray(addr.String())

	cnt := make(chan int, ip_concurrent)
	var wg sync.WaitGroup

	var alive_cnt = 0
	for i, ip := range addresses {
		cnt <- 1
		wg.Add(1)

		// fmt.Println("scanning:", ip)
		go func() {
			if noAliveDetect || isAlive(ip) {
				alive_cnt++
				if !noPortDetect {
					scanIp(ip, ch)
				}
			}
			wg.Done()
			<-cnt
		}()

		progress := float64(i+1) / float64(len(addresses)) * 100
		fmt.Printf("Progress: %.2f%%, Alive Count: %d\n", progress, alive_cnt)
	}

	wg.Wait() // wait for all goroutines to finish
	close(ch)
}

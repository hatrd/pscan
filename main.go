package main

import (
	"fmt"

	"github.com/hatrd/pscan/scan"
	"github.com/hatrd/pscan/ui"
)

func main() {
	ch := make(chan string)
	subnet := ui.GetSubNet()
	go scan.Scan(subnet, ch)
	for s := range ch {
		fmt.Println(s)
	}
}

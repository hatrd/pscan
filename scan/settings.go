package scan

import (
	"flag"
	"fmt"
	"time"
)

var ports = map[string]int{
	// "v2ray_local_socks": 10808, // 旧版不区分 local 和 lan，新版区分了
	// "v2ray_local_http":  10809,
	"v2ray_lan_socks": 10810,
	"v2ray_lan_http":  10811,
	"clash_all":       7890,
	"clash_socks":     7891,
	"clash_verge":     7897,
	"v2raya":          2017,
	"v2raya_socks":    20170,
	"v2raya_http":     20171,
}

var ip_concurrent = 25
var port_concurrent = 5
var output_file_name = fmt.Sprint("p_" + time.Now().Format("2006-01-02_15-04-05") + ".txt")

var noAliveDetect = false
var noPortDetect = false
var port_deadline = 3 * time.Second
var icmp_dealine = 2 * time.Second

func init() {
	flag.IntVar(&ip_concurrent, "i", ip_concurrent, "Number of concurrent IP scans")
	flag.IntVar(&port_concurrent, "p", port_concurrent, "Number of concurrent port scans")
	flag.StringVar(&output_file_name, "o", output_file_name, "Output file name")
	flag.BoolVar(&noAliveDetect, "na", noAliveDetect, "Disable alive detection")
	flag.BoolVar(&noPortDetect, "np", noPortDetect, "Disable port detection")
	flag.DurationVar(&port_deadline, "pd", port_deadline, "Port scan deadline")
	flag.DurationVar(&icmp_dealine, "id", icmp_dealine, "ICMP scan deadline")

	flag.Parse()
}

package main

import (
	"encoding/base64"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/lair-framework/go-nmap"
)

var (
	nmapFile    string
	outputFile  string
	format      string
	alive       bool
	port        int
	service     string
	nmapParser  nmap.NmapRun
	preHtml     string
	html        string
	csv         string
	postHtml    string
	info        string
	hostname    string
	ip          string
	outputHosts []outputHost
)

type outputHost struct {
	Ip       string
	Hostname string
	Protocol string
	PortId   int
	State    string
	Service  string
	Info     string
}

func printBanner() {
	bannerBase64 := `ICAgIF9fX18uICBfX19fXyAgICBfX19fX19fIF9fX19fX19fX18gDQogICAgfCAgICB8IC8gIF8gIFwgICBcICAgICAgXFxfX19fX18gICBcDQogICAgfCAgICB8LyAgL19cICBcICAvICAgfCAgIFx8ICAgICBfX18vDQovXF9ffCAgICAvICAgIHwgICAgXC8gICAgfCAgICBcICAgIHwgICAgDQpcX19fX19fX19cX19fX3xfXyAgL1xfX19ffF9fICAvX19fX3wgICAgDQogICAgICAgICAgICAgICAgIFwvICAgICAgICAgXC8gICAgICAgICAg`
	banner, _ := base64.StdEncoding.DecodeString(bannerBase64)
	fmt.Println(string(banner))
	fmt.Println("\tJust Another Nmap Parser")
	fmt.Println("Mattia Reggiani - https://github.com/mattiareggiani/janp - info@mattiareggiani.com")
	fmt.Println()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	printBanner()
	flag.StringVar(&nmapFile, "file", "", "Input Nmap XML file")
	flag.StringVar(&outputFile, "output", "output", "Output file")
	flag.StringVar(&format, "format", "", "Format of output file (HTML or CSV)")
	flag.BoolVar(&alive, "alive", false, "Print to stout the alive hosts")
	flag.StringVar(&service, "service", "", "Print to stout the hosts with the service open")
	flag.IntVar(&port, "port", 0, "Print to stout the hosts with the port open")

	flag.Parse()

	if !(len(nmapFile) > 0) {
		fmt.Println("[-] Insert the Nmap XML file")
		os.Exit(0)
	}
	xmlFile, err := os.Open(nmapFile)
	check(err)
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	err = xml.Unmarshal(byteValue, &nmapParser)
	if err != nil {
		fmt.Println("[-] Invalid Nmap XML file")
		panic(err)
	}
	fmt.Println("[*] Looks like a valid Nmap XML file")
	hosts := nmapParser.Hosts
	fmt.Println("[*] Parsing...")
	t0 := time.Now()
	for _, h := range hosts {
		hostname = ""
		ip = ""
		for _, n := range h.Hostnames {
			hostname += n.Name
		}
		for _, i := range h.Addresses {
			ip += i.Addr
		}
		for _, p := range h.Ports {
			info = fmt.Sprintf("%v %v %v", p.Service.Product, p.Service.Version, p.Service.ExtraInfo)
			outputHosts = append(outputHosts, outputHost{Ip: ip, Hostname: hostname, Protocol: p.Protocol, PortId: p.PortId, State: p.State.State, Service: p.Service.Name, Info: info})
		}
	}
	t1 := time.Now()
	fmt.Printf("[*] Parsed in %v\n", t1.Sub(t0))
	for _, h := range outputHosts {
		html += fmt.Sprintf("<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>", h.Ip, h.Hostname, h.Protocol, h.PortId, h.State, h.Service, h.Info)
		csv += fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,\n", h.Ip, h.Hostname, h.Protocol, h.PortId, h.State, h.Service, h.Info)

		if alive {
			// TODO
			fmt.Println("[-] Alive hosts is not implemented yet")
		}
		if (port != 0) && h.PortId == port {
			fmt.Println(h.Ip)
		}
		if len(service) > 0 && h.Service == service {
			fmt.Println(ip)
		}
	}

	if strings.EqualFold(format, "html") {
		fileName := outputFile + ".html"
		fH, err := os.Create(fileName)
		check(err)
		defer fH.Close()
		preHtml = `
<html>
<head>
<link rel="stylesheet" type="text/css" href="https://cdn.datatables.net/v/dt/jq-3.3.1/dt-1.10.18/af-2.3.0/kt-2.4.0/r-2.2.2/sc-1.5.0/sl-1.2.6/datatables.min.css"/>
 
<script type="text/javascript" src="https://cdn.datatables.net/v/dt/jq-3.3.1/dt-1.10.18/af-2.3.0/kt-2.4.0/r-2.2.2/sc-1.5.0/sl-1.2.6/datatables.min.js"></script>


</head>
<body>

<table id="example" class="display" style="width:100%">
        <thead>
            <tr>
                <th>Ip</th>
                <th>Hostname</th>
                <th>Protocol</th>
                <th>Port</th>
                <th>State</th>
                <th>Service</th>
                <th>Version</th>
            </tr>
        </thead>
        <tbody>
`
		postHtml += `
	</tbody>
    </table>
<script>
    $(document).ready(function() {
    $('#example').DataTable();
} );
</script>
</body>
</html>
	`
		b, err := fH.WriteString(preHtml + html + postHtml)
		fmt.Printf("[*] Wrote %d bytes\n", b)
		fmt.Printf("[+] HTML file successfully created: %v\n", fileName)
		fH.Sync()
	}
	if strings.EqualFold(format, "csv") {
		fileName := outputFile + ".csv"
		fC, err := os.Create(fileName)
		check(err)
		defer fC.Close()
		csv = "Ip,Hostname,Protocol,Port,State,Service,Info,\n" + csv
		b, err := fC.WriteString(csv)
		fmt.Printf("[*] Wrote %d bytes\n", b)
		fmt.Printf("[+] CSV file successfully created: %v\n", fileName)
		fC.Sync()
	}
}

package findport

import (
	"github.com/Ullaakut/nmap"
	"strconv"
	"strings"
	"time"
)

func (r *mas_scan) getsevice(ip string) error {
	retry := true
again:
	scanner, err := nmap.NewScanner(
		nmap.WithCustomArguments("--disable-arp-ping"),
		nmap.WithCustomArguments("-sV"),
		nmap.WithCustomArguments("--version-intensity"),
		nmap.WithCustomArguments(r.nmapsv), // 0-9
		nmap.WithCustomArguments("-Pn"),
		nmap.WithCustomArguments("--script-args"),
		nmap.WithCustomArguments("-p"),
		nmap.WithCustomArguments(r.masscan_port),
		nmap.WithCustomArguments("-n"),
		nmap.WithCustomArguments("-r"),
		nmap.WithCustomArguments(ip),
	)
	if err != nil {
		return err
	}
	run, _, err := scanner.Run()
	if err != nil {
		if retry && strings.Contains(err.Error(), "unable to parse nmap output") {
			retry = false
			time.Sleep(2 * time.Second)
			goto again
		}
		return err
	}
	for _, hosts := range run.Hosts {
		for _, ports := range hosts.Ports {
			if len(ports.Service.ServiceFP) != 0 && ports.Service.Name != "unknown" {
				ports.Service.Name = ports.Service.Name + "?"
			}
			if len(ports.Service.Tunnel) != 0 {
				ports.Service.Name = ports.Service.Tunnel + "/" + ports.Service.Name
			}
			if strings.Contains(ports.Service.ServiceFP, "HTTP") {
				ports.Service.Name = "http"
			}
			r.Services = append(r.Services, Service{IP: ip, Port: strconv.Itoa(int(ports.ID)), Protocol: ports.Service.Name, Version: strings.Trim(ports.Service.Product+" "+ports.Service.Version, " "), ServiceFP: ports.Service.ServiceFP})
		}
	}
	return nil
}

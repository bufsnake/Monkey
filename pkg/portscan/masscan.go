package findport

import (
	"errors"
	"github.com/bufsnake/Monkey/config"
	"github.com/bufsnake/go-masscan"
	"strconv"
	"strings"
)

type mas_scan struct {
	conf         config.Config
	port         string
	nmapsv       string
	masscan_rate string
	masscan_port string
	Services     []Service
}

func (f *mas_scan) PortScan(ip string) error {
	if !f.conf.NoPing {
		ping, err := pingscan(ip)
		if err != nil {
			return err
		}
		if !ping {
			return errors.New("ping scan " + ip + " not alive")
		}
	}
	err := f.getport(ip)
	if err != nil {
		return err
	}
	err = f.getsevice(ip)
	if err != nil {
		return err
	}
	return nil
}

func (f *mas_scan) GetService() []Service {
	return f.Services
}

func (r *mas_scan) getport(ip string) error {
	m := masscan.New()
	m.SetPorts(r.port)
	m.SetRate(r.masscan_rate)
	m.SetArgs([]string{ip, "--wait", "1"}...)
	err := m.Run()
	if err != nil {
		return err
	}
	results, err := m.Parse()
	if err != nil {
		return err
	}
	count := 0
	for _, result := range results {
		for _, ip := range result.Ports {
			r.masscan_port += ip.Portid + ","
			count++
		}
	}
	if count > 1000 {
		return errors.New("too manay ports, there are " + strconv.Itoa(count) + " ports in total.")
	}
	if r.masscan_port == "" {
		return errors.New("not found port")
	}
	r.masscan_port = strings.Trim(r.masscan_port, ",")
	return nil
}

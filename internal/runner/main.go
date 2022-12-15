package runner

import (
	"fmt"
	"github.com/bufsnake/Monkey/config"
	"github.com/bufsnake/Monkey/pkg/blasting"
	"github.com/bufsnake/Monkey/pkg/log"
	"github.com/bufsnake/Monkey/pkg/parseip"
	"github.com/bufsnake/Monkey/pkg/portscan"
	"github.com/bufsnake/Monkey/pkg/sitescan"
	"strings"
	"sync"
	"time"
)

type Runner struct {
	conf     config.Config
	probes   [][2]uint32
	services map[string][]portscan.Service
}

func NewRunner(conf config.Config, probes [][2]uint32) *Runner {
	return &Runner{conf: conf, probes: probes, services: make(map[string][]portscan.Service)}
}

func (c *Runner) Run() {
	go log.Bar()
	
	// 端口扫描
	portwait := sync.WaitGroup{}
	portchan := make(chan string, c.conf.Thread)
	for i := 0; i < c.conf.Thread; i++ {
		portwait.Add(1)
		go c.portscan(&portwait, portchan)
	}
	m := make(map[uint32]bool)
	for i := 0; i < len(c.probes); i++ {
		for j := c.probes[i][0]; j <= c.probes[i][1]; j++ {
			if _, ok := m[j]; ok {
				continue
			}
			portchan <- parseip.UInt32ToIP(j)
		}
	}
	close(portchan)
	portwait.Wait()
	
	log.STOP = true
	time.Sleep(time.Second)
	
	fmt.Println()
	// HTTP服务探测
	fmt.Println("HTTP:")
	var n_weak = []string{"SSH", "SSL/SSH", "FTP", "SSL/FTP", "MYSQL", "NAGIOS-NSCA", "SSL/MYSQL", "SSL/NAGIOS-NSCA", "MS-SQL-S", "SSL/MS-SQL-S", "REDIS", "SSL/REDIS", "MONGOD", "SSL/MONGOD", "MONGODB", "SSL/MONGODB", "POSTGRESQL", "SSL/POSTGRESQL", "MICROSOFT-DS", "SSL/MICROSOFT-DS", "SNMP", "SSL/SNMP", "DOCKER", "SSL/DOCKER", "ZOOKEEPER", "SSL/ZOOKEEPER", "MEMCACHE", "SSL/MEMCACHE", "MS-WBT-SERVER", "SSL/MS-WBT-SERVER", "TELNET", "SSL/TELNET", "ORACLE", "SSL/ORACLE", "DOCKER", "ZOOKEEPER", "MEMCACHED"}
	for ip, services := range c.services {
		for i := 0; i < len(services); i++ {
			if c.conf.Blast && c.exist(n_weak, strings.ToUpper(strings.Trim(services[i].Protocol, "?")), true) {
				continue
			}
			// HTTP 探测
			httpx := sitescan.NewHttpx(ip+":"+services[i].Port, c.conf.Timeout)
			err := httpx.Run()
			if err != nil {
				continue
			}
			for j := 0; j < len(httpx.URLS); j++ {
				fmt.Println(httpx.URLS[j].GetUrl())
			}
		}
	}
	if c.conf.Blast {
		fmt.Println("WEAK PASSWORD:")
	}
	// 弱口令扫描
	for ip, services := range c.services {
		for i := 0; i < len(services); i++ {
			if !c.conf.Blast || c.exist(n_weak, strings.ToUpper(strings.Trim(services[i].Protocol, "?")), true) {
				continue
			}
			// 弱口令探测
			blast, err := blasting.NewBlast(strings.Trim(services[i].Protocol, "?"), ip, services[i].Port)
			if err != nil {
				continue
			}
			weakurl, err := blast.Connect()
			if err != nil {
				continue
			}
			fmt.Println(weakurl)
		}
	}
	// 全部端口
	mx := make(map[string]bool)
	ports := make([]string, 0)
	for _, services := range c.services {
		for i := 0; i < len(services); i++ {
			if _, ok := mx[services[i].Port]; ok {
				continue
			}
			mx[services[i].Port] = true
			ports = append(ports, services[i].Port)
		}
	}
	if len(ports) == 0 {
		return
	}
	fmt.Print("PORTS:\n" + strings.Join(ports, ","))
}

// strong 强判断
func (c *Runner) exist(arr []string, data string, strong bool) bool {
	for i := 0; i < len(arr); i++ {
		if strong {
			if arr[i] == data {
				return true
			}
		} else if strings.Contains(data, arr[i]) {
			return true
		}
	}
	return false
}

func (c *Runner) portscan(wait *sync.WaitGroup, portchan chan string) {
	defer wait.Done()
	for ip := range portchan {
		scan := portscan.NewPortScan(c.conf)
		err := scan.PortScan(ip)
		if err != nil {
			log.Println(ip, err.Error()+strings.Repeat(" ", 18))
			continue
		}
		services := scan.GetService()
		log.UpdateTotalCount(len(services))
		for i := 0; i < len(services); i++ {
			port, protocol, version := services[i].Port, services[i].Protocol, services[i].Version
			if strings.TrimSpace(services[i].ServiceFP) != "" {
				protocol = strings.Trim(protocol, "?") + "?"
			}
			log.Println(ip, port, protocol, version)
		}
		if _, ok := c.services[ip]; !ok {
			c.services[ip] = make([]portscan.Service, 0)
		}
		c.services[ip] = append(c.services[ip], services...)
	}
}

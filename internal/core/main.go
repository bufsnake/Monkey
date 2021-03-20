package core

import (
	"github.com/bufsnake/Monkey/config"
	"github.com/bufsnake/Monkey/pkg/log"
	findport "github.com/bufsnake/Monkey/pkg/portscan"
	web "github.com/bufsnake/Monkey/pkg/sitescan"
	"strconv"
	"strings"
	"sync"
)

type core struct {
	conf  config.Config
	allip []string
}

func NewCore(conf config.Config, allip []string) core {
	return core{conf: conf, allip: allip}
}

func (c *core) Start() {
	go log.Bar()
	portwait := sync.WaitGroup{}
	httpwait := sync.WaitGroup{}
	portchan := make(chan string, c.conf.Thread)
	httpchan := make(chan findport.Service, c.conf.Thread*2)
	for i := 0; i < c.conf.Thread; i++ {
		portwait.Add(1)
		httpwait.Add(1)
		go c.portscan(&portwait, portchan, httpchan)
		go c.httpchan(&httpwait, httpchan)
	}
	for i := 0; i < len(c.allip); i++ {
		portchan <- c.allip[i]
	}
	close(portchan)
	portwait.Wait()
	close(httpchan)
	httpwait.Wait()
}

var httpx = []string{"sun-answerbook", "http", "tcpwrapped", "caldav", "sip", "rtsp", "soap"}

// strong 强判断
func exist(arr []string, data string, strong bool) bool {
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

func (c *core) portscan(portwait *sync.WaitGroup, portchan chan string, httpchan chan findport.Service) {
	defer portwait.Done()
	for ip := range portchan {
		scan, err := findport.NewPortScan(c.conf)
		if err != nil {
			log.Println(err)
			continue
		}
		err = scan.PortScan(ip)
		if err != nil {
			log.Println(ip, err.Error()+strings.Repeat(" ", 18))
			continue
		}
		log.UpdateTotalCount(len(scan.Services))
		for i := 0; i < len(scan.Services); i++ {
			if exist(httpx, scan.Services[i].Protocol, false) { // 判断是否为http
				httpchan <- scan.Services[i]
			} else { // 其他情况直接打印
				log.Println(ip, scan.Services[i].Port, scan.Services[i].Protocol, scan.Services[i].Version, scan.Services[i].ServiceFP)
			}
		}
	}
}

func (c *core) httpchan(httpwait *sync.WaitGroup, httpchan chan findport.Service) {
	defer httpwait.Done()
	for httpip := range httpchan {
		log.Println(httpip.IP, httpip.Port, httpip.Protocol, httpip.Version, httpip.ServiceFP)
		newHttpx := web.NewHttpx(httpip.IP+":"+httpip.Port, c.conf.Timeout)
		err := newHttpx.Run()
		if err != nil {
			continue
		}
		for i := 0; i < len(newHttpx.URLS); i++ {
			config.AppendUrls([]string{strconv.Itoa(newHttpx.URLS[i].GetCode()), newHttpx.URLS[i].GetUrl(), newHttpx.URLS[i].GetTitle(), strconv.Itoa(newHttpx.URLS[i].GetLength())})
		}
	}
}

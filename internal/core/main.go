package core

import (
	"github.com/bufsnake/Monkey/config"
	"github.com/bufsnake/Monkey/pkg/blasting"
	"github.com/bufsnake/Monkey/pkg/log"
	findport "github.com/bufsnake/Monkey/pkg/portscan"
	web "github.com/bufsnake/Monkey/pkg/sitescan"
	"net"
	"strconv"
	"strings"
	"sync"
)

type core struct {
	conf  config.Config
	allip [][2]uint32
}

func NewCore(conf config.Config, allip [][2]uint32) core {
	return core{conf: conf, allip: allip}
}

func (c *core) Start() {
	go log.Bar()
	portwait := sync.WaitGroup{}
	httpwait := sync.WaitGroup{}
	blastwait := sync.WaitGroup{}
	portchan := make(chan string, c.conf.Thread)
	httpchan := make(chan findport.Service, c.conf.Thread*2)
	blastchan := make(chan findport.Service, c.conf.Thread*2)
	for i := 0; i < c.conf.Thread; i++ {
		portwait.Add(1)
		httpwait.Add(1)
		blastwait.Add(1)
		go c.portscan(&portwait, portchan, httpchan, blastchan)
		go c.httpchan(&httpwait, httpchan)
		go c.blastchan(&blastwait, blastchan)
	}
	for i := 0; i < len(c.allip); i++ {
		for j := c.allip[i][0]; j < c.allip[i][1]; j++ {
			portchan <- UInt32ToIP(j)
		}
	}
	close(portchan)
	portwait.Wait()
	close(httpchan)
	close(blastchan)
	httpwait.Wait()
	blastwait.Wait()
}

var httpx = []string{"sun-answerbook", "http", "tcpwrapped", "caldav", "sip", "rtsp", "soap"}
var blast = []string{"docker", "zookeeper", "memcached", "ssh", "ftp", "mysql", "nagios-nsca", "ms-sql-s", "redis", "mongodb", "postgresql", "microsoft-ds", "ssl/ssh", "ssl/ftp", "ssl/mysql", "ssl/nagios-nsca", "ssl/ms-sql-s", "ssl/redis", "ssl/mongodb", "ssl/postgresql", "ssl/microsoft-ds"}

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

func (c *core) portscan(portwait *sync.WaitGroup, portchan chan string, httpchan, blastchan chan findport.Service) {
	defer portwait.Done()
	for ip := range portchan {
		scan := findport.NewPortScanX(c.conf)
		err := scan.PortScan(ip)
		if err != nil {
			log.Println(ip, err.Error()+strings.Repeat(" ", 18))
			continue
		}
		log.UpdateTotalCount(len(scan.GetService()))
		for i := 0; i < len(scan.GetService()); i++ {
			if exist(httpx, scan.GetService()[i].Protocol, false) { // 判断是否为http
				httpchan <- scan.GetService()[i]
			} else if c.conf.Blast && exist(blast, scan.GetService()[i].Protocol, true) { // 判断是否为blast
				blastchan <- scan.GetService()[i]
			} else { // 其他情况直接打印
				log.Println(ip, scan.GetService()[i].Port, scan.GetService()[i].Protocol, scan.GetService()[i].Version, scan.GetService()[i].ServiceFP)
			}
			config.AppendAllPort(scan.GetService()[i].Port)
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

func (c *core) blastchan(blastwait *sync.WaitGroup, blastchan chan findport.Service) {
	defer blastwait.Done()
	for blastx := range blastchan {
		log.Println(blastx.IP, blastx.Port, blastx.Protocol, blastx.Version, blastx.ServiceFP)
		newBlast, err := blasting.NewBlast(blastx.Protocol, blastx.IP, blastx.Port)
		if err != nil {
			continue
		}
		connect, err := newBlast.Connect()
		if err != nil {
			continue
		}
		config.AppendBlast(connect)
	}
}

func IP2UInt32(ipnr string) uint32 {
	bits := strings.Split(ipnr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32
	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)
	return sum
}

func UInt32ToIP(intIP uint32) string {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

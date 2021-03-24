package findport

import (
	"errors"
	"fmt"
	"github.com/bufsnake/Monkey/config"
	"net"
	"strconv"
	"sync"
	"time"
)

type socket struct {
	conf         config.Config
	port         string
	nmapsv       string
	masscan_rate string
	masscan_port string
	Services     []Service

	ip string
}

func (f *socket) GetService() []Service {
	return f.Services
}

func (s *socket) PortScan(ip string) error {
	if !s.conf.NoPing {
		ping, err := pingscan(ip)
		if err != nil {
			return err
		}
		if !ping {
			return errors.New("ping scan " + ip + " not alive")
		}
	}
	err := s.getport(ip)
	if err != nil {
		return err
	}
	temp := mas_scan{
		port:         s.port,
		nmapsv:       s.nmapsv,
		Services:     s.Services,
		masscan_port: s.masscan_port,
	}
	err = temp.getsevice(ip)
	if err != nil {
		return err
	}
	s.Services = temp.Services
	return nil
}

func (r *socket) getport(ip string) error {
	r.ip = ip
	wait := sync.WaitGroup{}
	ports := make(chan int, 200)
	port_result := make(chan int, 1000)
	for i := 0; i < r.conf.SocketThreads; i++ {
		wait.Add(1)
		go r.connect(&wait, ports, port_result)
	}
	pw := sync.WaitGroup{}
	pw.Add(1)

	go func() {
		pw.Done()
		for port := range port_result {
			r.masscan_port += strconv.Itoa(port) + ","
		}
	}()

	for i := 1; i < 65536; i++ {
		ports <- i
	}

	close(ports)
	if r.masscan_port == "" {
		return errors.New("not found port")
	}
	wait.Wait()
	close(port_result)
	pw.Wait()
	return nil
}

func (r *socket) connect(wait *sync.WaitGroup, ports, port_result chan int) {
	defer wait.Done()
	for port := range ports {
		con, err := net.DialTimeout("tcp4", fmt.Sprintf("%s:%d", r.ip, port), time.Duration(1)*time.Second)
		if err == nil {
			_ = con.Close()
			port_result <- port
		}
	}
}

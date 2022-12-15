package main

import (
	"flag"
	"fmt"
	"github.com/bufsnake/Monkey/config"
	"github.com/bufsnake/Monkey/internal/runner"
	"github.com/bufsnake/Monkey/pkg/log"
	"github.com/bufsnake/Monkey/pkg/parseip"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 3 / 4)
	var rLimit syscall.Rlimit
	rLimit.Max = 999999
	rLimit.Cur = 999999
	if runtime.GOOS == "darwin" {
		rLimit.Cur = 10240
	}
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_ = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	localip := []string{}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localip = append(localip, ipnet.IP.String())
			}
		}
	}
	fmt.Println("Linux: sudo sysctl -w net.ipv4.ping_group_range=\"0 2147483647\"")
	fmt.Println("Local IP is", localip, ",ulimit -n is", rLimit.Cur)
	fmt.Println("内网扫描需要配置网关，网关IP为任意一台存活IP即可")
	conf := config.Config{}
	flag.StringVar(&conf.Target, "target", "", "specify target ip")
	flag.StringVar(&conf.Targets, "targets", "", "specify target ip from file")
	flag.IntVar(&conf.Thread, "thread", 6, "specify scan threads")
	flag.IntVar(&conf.Timeout, "timeout", 10, "web request timeout")
	flag.BoolVar(&conf.Blast, "blast", false, "blast service")
	flag.StringVar(&conf.NmapSV, "nmap-sv", "2", "nmap version intensity,optional 0~9")
	flag.StringVar(&conf.Port, "port", "0-65535", "specify scan ports")
	flag.StringVar(&conf.MasScanRate, "masscan-rate", "1000", "masscan rate")
	flag.Parse()
	if conf.Target == "" && conf.Targets == "" {
		flag.Usage()
		return
	}
	probes := make([][2]uint32, 0)
	if conf.Target != "" {
		start, end, err := parseip.ParseIP(conf.Target)
		if err != nil {
			fmt.Println(conf.Target, err)
			return
		}
		probes = append(probes, [2]uint32{start, end})
	} else if conf.Targets != "" {
		file, err := os.ReadFile(conf.Targets)
		if err != nil {
			log.Println(err)
			return
		}
		split := strings.Split(string(file), "\n")
		for i := 0; i < len(split); i++ {
			split[i] = strings.Trim(split[i], "\r")
			if len(split[i]) == 0 {
				continue
			}
			start, end, err := parseip.ParseIP(split[i])
			if err != nil {
				log.Println(err)
				return
			}
			probes = append(probes, [2]uint32{start, end})
		}
	}
	
	var count uint32
	for i := 0; i < len(probes); i++ {
		count += probes[i][1] - probes[i][0] + 1
	}
	log.SetTotalCount(int(count))
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		select {
		case _ = <-c:
			fmt.Printf("\033[?25h")
			os.Exit(1)
		}
	}()
	runner.NewRunner(conf, probes).Run()
	fmt.Println()
	fmt.Printf("\033[?25h")
}

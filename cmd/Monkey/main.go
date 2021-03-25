package main

import (
	"flag"
	"fmt"
	"github.com/bufsnake/Monkey/config"
	"github.com/bufsnake/Monkey/internal/core"
	"github.com/bufsnake/Monkey/pkg/log"
	"github.com/bufsnake/Monkey/pkg/parseip"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strings"
	"syscall"
)

func main() {
	// 开启多核模式
	runtime.GOMAXPROCS(runtime.NumCPU() * 3 / 4)
	// 关闭 GIN Debug模式
	// 设置工具可打开的文件描述符
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
	fmt.Println("Local IP is", localip, ",ulimit -n is", rLimit.Cur)

	conf := config.Config{}
	flag.StringVar(&conf.IP, "t", "", "specify target ip")
	flag.StringVar(&conf.File, "f", "", "specify target ip from file")
	flag.IntVar(&conf.Thread, "r", 6, "specify scan rate")
	flag.IntVar(&conf.Timeout, "w", 10, "web request timeout")
	flag.IntVar(&conf.SocketThreads, "st", 200, "socket port scan thread")
	flag.BoolVar(&conf.Blast, "b", false, "whether to blast")
	flag.BoolVar(&conf.NoPing, "np", false, "whether to ping scan ip alive")
	flag.BoolVar(&conf.Masscan, "um", false, "use masscan(IM: use in public network)")
	flag.StringVar(&conf.NmapsV, "n", "2", "nmap version intensity,optional 0~9")
	flag.StringVar(&conf.Port, "p", "0-65535", "specify scan ports")
	flag.StringVar(&conf.MasscanRate, "m", "1000", "masscan rate")
	flag.Parse()
	if conf.IP == "" && conf.File == "" {
		fmt.Println("Linux: sudo sysctl -w net.ipv4.ping_group_range=\"0 2147483647\"")
		flag.Usage()
		return
	}
	allip := [][2]uint32{}
	if conf.IP != "" {
		startx, endx, err := parseip.ParseIP(strings.Trim(conf.IP, " "))
		if err != nil {
			log.Println(err)
			return
		}
		allip = append(allip, [2]uint32{startx, endx})
	} else if conf.File != "" {
		file, err := ioutil.ReadFile(conf.File)
		if err != nil {
			log.Println(err)
			return
		}
		split := strings.Split(string(file), "\n")
		for i := 0; i < len(split); i++ {
			trim := strings.Trim(split[i], "\r")
			if len(trim) == 0 {
				continue
			}
			if len(trim) >= 7 {
				startx, endx, err := parseip.ParseIP(trim)
				if err != nil {
					log.Println(err)
					return
				}
				allip = append(allip, [2]uint32{startx, endx})
			} else {
				log.Println("error data", trim)
				return
			}
		}
	}

	var count uint32
	for i := 0; i < len(allip); i++ {
		count += allip[i][1] - allip[i][0] + 1
	}
	log.SetOutputCSV(conf.Output)
	log.SetTotalCount(int(count))

	xcore := core.NewCore(conf, allip)
	xcore.Start()
	fmt.Println()
	for i := 0; i < len(config.Urldatas); i++ {
		if i == 0 {
			fmt.Println("Urls Results:")
		}
		fmt.Println("["+aurora.BrightMagenta(config.Urldatas[i][0]).String()+"]", "["+aurora.BrightWhite(config.Urldatas[i][1]).String()+"]", "["+aurora.BrightGreen(config.Urldatas[i][3]).String()+"]", "["+aurora.BrightCyan(config.Urldatas[i][2]).String()+"]")

	}
	for i := 0; i < len(config.Blastdatas); i++ {
		if i == 0 {
			fmt.Println("Blasting Results:")
		}
		fmt.Println(config.Blastdatas[i])
	}
	for i := 0; i < len(config.AllPort); i++ {
		if i == 0 {
			fmt.Println("AllPort Results:")
		}
		fmt.Print(config.AllPort[i])
		if i+1 == len(config.AllPort) {
			break
		}
		fmt.Print(",")
	}
	fmt.Println()
}

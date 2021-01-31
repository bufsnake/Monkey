package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/Ullaakut/nmap"
	"github.com/bufsnake/go-masscan"
	"github.com/bufsnake/parseip"
	"github.com/fatih/color"
	_ "github.com/google/gopacket"
	_ "github.com/google/gopacket/layers"
	"github.com/kyokomi/emoji"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var (
	ips        string
	thread     int
	alive      []string
	alive_lock sync.Mutex
	version    string
	portthread int
	masscaner  bool
	alllink    []string
	file       string
	nmapalive  bool
	lock       sync.Mutex
)

func init() {
	var rLimit syscall.Rlimit
	rLimit.Max = 999999
	rLimit.Cur = 999999
	if runtime.GOOS == "darwin" {
		rLimit.Cur = 10240
	}
	err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		color.Red("Error Setting ulimit " + err.Error())
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		color.Red("Error Getting ulimit " + err.Error())
	}
	flag.StringVar(&ips, "i", "", "指定IP")
	flag.IntVar(&thread, "t", 50, "指定线程,默认50")
	flag.IntVar(&portthread, "p", 5, "指定端口扫描线程,默认50")
	flag.StringVar(&version, "v", "2", "指定-sV详细程度0-9")
	flag.BoolVar(&masscaner, "m", false, "指定是否使用masscan进行端口扫描")
	flag.StringVar(&file, "f", "", "从文件中获取IP")
	flag.BoolVar(&nmapalive, "nmap-alive", false, "是否使用nmap进行探活")
	flag.Parse()
	if ips == "" && file == "" {
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	emoji.Println(":new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon::new_moon:")
	start := time.Now()
	ip := []string{}
	if ips != "" {
		ip = parseip.ParseIP(ips)
	} else if file != "" {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			color.Red(err.Error())
			os.Exit(1)
		}
		all_ip := strings.Split(string(data), "\n")
		for _, val := range all_ip {
			if val == "" {
				continue
			}
			temp := parseip.ParseIP(val)
			for _, temp_val := range temp {
				if !Exist(ip, temp_val) {
					ip = append(ip, temp_val)
				}
			}
		}
	}
	emoji.Println(":beer: total " + strconv.Itoa(len(ip)) + " ip\n:beer: time " + (time.Now().Sub(start)).String())
	emoji.Println(":waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon:")
	color.Yellow("Start to scan alive")
	start = time.Now()
	emoji.Println(":waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon::waning_crescent_moon:")
	Scan(ip, len(ip), 50)
	emoji.Println(":last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon:")
	emoji.Println(":beer: " + strconv.Itoa(len(alive)) + " alive\n:beer: time " + (time.Now().Sub(start)).String())
	color.Yellow("Start to scan service")
	now := time.Now()
	emoji.Println(":last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon::last_quarter_moon:")
	ServiceScan()
	emoji.Println(":waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon:")
	emoji.Println(":beer: " + time.Now().Sub(now).String())
	color.Yellow("Start to scan web")
	now = time.Now()
	emoji.Println(":waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon::waning_gibbous_moon:")
	WebScan()
	emoji.Println(":full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon::full_moon:")
	emoji.Println(":beer: " + time.Now().Sub(now).String())
}

func WebScan() {
	var WebWait sync.WaitGroup
	ipchan := make(chan string, 30000)
	iplistlen := len(alllink)
	DistributionIPWait := make(chan int)
	for i := 0; i < iplistlen; i += thread {
		if i+thread > iplistlen {
			go DistributionIP(alllink[i:iplistlen], ipchan, DistributionIPWait)
			break
		}
		go DistributionIP(alllink[i:i+thread], ipchan, DistributionIPWait)
	}
	for i := 0; i < thread; i++ {
		WebWait.Add(1)
		go GetWeb(&WebWait, ipchan)
	}
	for i := 0; i < iplistlen; i += thread {
		<-DistributionIPWait
	}
	close(ipchan)
	WebWait.Wait()
}

func Scan(iplist []string, iplistlen, threads int) {
	var SurviveWait sync.WaitGroup
	ipchan := make(chan string, 30000)
	DistributionIPWait := make(chan int)
	for i := 0; i < iplistlen; i += threads {
		if i+threads > iplistlen {
			go DistributionIP(iplist[i:iplistlen], ipchan, DistributionIPWait)
			break
		}
		go DistributionIP(iplist[i:i+threads], ipchan, DistributionIPWait)
	}
	for i := 0; i < threads; i++ {
		SurviveWait.Add(1)
		go Survive(&SurviveWait, ipchan)
	}
	for i := 0; i < iplistlen; i += threads {
		<-DistributionIPWait
	}
	close(ipchan)
	SurviveWait.Wait()
}

func ServiceScan() {
	iplist := alive
	iplistlen := len(iplist)
	ipchan := make(chan string, thread)
	DistributionIPWait := make(chan int)
	for i := 0; i < iplistlen; i += thread {
		if i+thread > iplistlen {
			go SDistributionIP(iplist[i:iplistlen], ipchan, DistributionIPWait)
			break
		}
		go SDistributionIP(iplist[i:i+thread], ipchan, DistributionIPWait)
	}

	var SurviveWait sync.WaitGroup
	for i := 0; i < thread; i++ {
		SurviveWait.Add(1)
		go TCPScan(&SurviveWait, ipchan)
	}
	for i := 0; i < iplistlen; i += thread {
		<-DistributionIPWait
	}
	close(ipchan)
	SurviveWait.Wait()
}

func SDistributionIP(iplist []string, ipchan chan string, finishflag chan int) {
	for i := 0; i < len(iplist); i++ {
		ipchan <- iplist[i]
	}
	finishflag <- 1
}

func DistributionIP(iplist []string, ipchan chan string, finishflag chan int) {
	for i := 0; i < len(iplist); i++ {
		ipchan <- iplist[i]
	}
	finishflag <- 1
}

func Survive(wait *sync.WaitGroup, ipchan chan string) {
	defer wait.Done()
	for ip := range ipchan {
		if !nmapalive {
			alive_lock.Lock()
			alive = append(alive, ip)
			alive_lock.Unlock()
		} else if SNSurvive(ip) == 1 {
			alive_lock.Lock()
			color.Green(ip)
			alive = append(alive, ip)
			alive_lock.Unlock()
		} else if PingSurvive(ip) == 1 {
			alive_lock.Lock()
			color.Green(ip)
			alive = append(alive, ip)
			alive_lock.Unlock()
		} else if ICMPSurvive(ip) == 1 {
			alive_lock.Lock()
			alive = append(alive, ip)
			color.Green(ip)
			alive_lock.Unlock()
		} else if ARPSurvive(ip) == 1 {
			alive_lock.Lock()
			alive = append(alive, ip)
			color.Green(ip)
			alive_lock.Unlock()
		} else if PMSurvive(ip) == 1 {
			alive_lock.Lock()
			alive = append(alive, ip)
			color.Green(ip)
			alive_lock.Unlock()
		}
	}
}

func SNSurvive(ip string) int {
	scanner, err := nmap.NewScanner(
		nmap.WithCustomArguments("-sn"),
		nmap.WithCustomArguments("-n"),
		nmap.WithCustomArguments(ip),
	)
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	run, _, err := scanner.Run()
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	if len(run.Hosts) != 0 {
		if run.Hosts[0].Status.State == "up" {
			return 1
		}
	}
	return 0
}

func PingSurvive(ip string) int {
	scanner, err := nmap.NewScanner(
		nmap.WithCustomArguments("-sn"),
		nmap.WithCustomArguments("-n"),
		nmap.WithCustomArguments("-PP"),
		nmap.WithCustomArguments("--disable-arp-ping"),
		nmap.WithCustomArguments(ip),
	)
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	run, _, err := scanner.Run()
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	if len(run.Hosts) != 0 {
		if run.Hosts[0].Status.State == "up" {
			return 1
		}
	}
	return 0
}

func ICMPSurvive(ip string) int {
	scanner, err := nmap.NewScanner(
		nmap.WithCustomArguments("--disable-arp-ping"),
		nmap.WithCustomArguments("-n"),
		nmap.WithCustomArguments("-sn"),
		nmap.WithCustomArguments("-PE"),
		nmap.WithCustomArguments(ip),
	)
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	run, _, err := scanner.Run()
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	if len(run.Hosts) != 0 {
		if run.Hosts[0].Status.State == "up" {
			return 1
		}
	}
	return 0
}

func ARPSurvive(ip string) int {
	scanner, err := nmap.NewScanner(
		nmap.WithCustomArguments("--disable-arp-ping"),
		nmap.WithCustomArguments("-n"),
		nmap.WithCustomArguments("-PR"),
		nmap.WithCustomArguments("-sn"),
		nmap.WithCustomArguments(ip),
	)
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	run, _, err := scanner.Run()
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	if len(run.Hosts) != 0 {
		if run.Hosts[0].Status.State == "up" {
			return 1
		}
	}
	return 0
}

func PMSurvive(ip string) int {
	scanner, err := nmap.NewScanner(
		nmap.WithCustomArguments("--disable-arp-ping"),
		nmap.WithCustomArguments("-n"),
		nmap.WithCustomArguments("-PM"),
		nmap.WithCustomArguments("-sn"),
		nmap.WithCustomArguments(ip),
	)
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	run, _, err := scanner.Run()
	if err != nil {
		color.Red(ip + " " + err.Error())
		return 0
	}
	if len(run.Hosts) != 0 {
		if run.Hosts[0].Status.State == "up" {
			return 1
		}
	}
	return 0
}

func TCPScan(wait *sync.WaitGroup, ipchan chan string) {
	defer wait.Done()
	for ip := range ipchan {
		temp := ""
		if !masscaner {
			port := TPortScan(ip)
			for i := 0; i < len(port); i++ {
				temp += port[i] + ","
			}
		} else {
			temp = MPortScan(ip)
		}
		if temp == "" {
			color.Blue(ip + " is alive,no port open")
			continue
		}
		scanner, err := nmap.NewScanner(nmap.WithCustomArguments("--disable-arp-ping"), nmap.WithCustomArguments("-sV"), nmap.WithCustomArguments("--version-intensity"), nmap.WithCustomArguments(version), nmap.WithCustomArguments("-Pn"), nmap.WithCustomArguments("-p"), nmap.WithCustomArguments(temp), nmap.WithCustomArguments("-n"), nmap.WithCustomArguments(ip))
		if err != nil {
			color.Red("nmap create error " + err.Error())
			os.Exit(1)
		}
		run, _, err := scanner.Run()
		if err != nil {
			color.Red(ip + " is open " + strings.Trim(temp, ",") + " port " + err.Error())
			continue
		}
		for _, hosts := range run.Hosts {
			for _, ports := range hosts.Ports {
				if len(ports.Service.Tunnel) != 0 {
					ports.Service.Name = ports.Service.Tunnel + "/" + ports.Service.Name
				}
				if len(ports.Service.ServiceFP) != 0 && ports.Service.Name != "unknown" {
					ports.Service.Name = ports.Service.Name + "?"
				}
				if strings.Contains(ports.Service.ServiceFP, "HTTP") {
					ports.Service.Name = "http"
				}
				link := ""
				if strings.Contains(ports.Service.Name, "http") || strings.Contains(ports.Service.Name, "tcpwrapped") || strings.Contains(ports.Service.Name, "caldav") || strings.Contains(ports.Service.Name, "sip") || strings.Contains(ports.Service.Name, "rtsp") || strings.Contains(ports.Service.Name, "soap") {
					link := ">>> WEB <<<"
					lock.Lock()
					alllink = append(alllink, ip+":"+strconv.Itoa(int(ports.ID)))
					lock.Unlock()
					fmt.Println(version, fmt.Sprintf("%-4s", ports.Protocol), fmt.Sprintf("%-15s", ip), fmt.Sprintf("%-5d", ports.ID), fmt.Sprintf("%-25s", ports.Service.Name), fmt.Sprintf("%-70s", strings.Trim(ports.Service.Product+" "+ports.Service.Version, " ")), fmt.Sprintf("%s", link))
				} else {
					if ports.Service.Name == "x11" {
						link = "xwd -root -screen -silent -display x:0 > o.xwd && convert o.xwd o.png"
					}
					fmt.Println(version, fmt.Sprintf("%-4s", ports.Protocol), fmt.Sprintf("%-15s", ip), fmt.Sprintf("%-5d", ports.ID), fmt.Sprintf("%-25s", ports.Service.Name), fmt.Sprintf("%-70s", strings.Trim(ports.Service.Product+" "+ports.Service.Version, " ")), fmt.Sprintf("%s", link))
				}
			}
		}
	}
}

func TPortScan(ip string) []string {
	ret := []string{}
	port_thread := portthread
	runtime.GOMAXPROCS(runtime.NumCPU() / 4 * 3)
	var wait = sync.WaitGroup{}
	ports := make(chan string, 60000)
	for i := 0; i < port_thread; i++ {
		wait.Add(1)
		go func(ports chan string, ip string) {
			defer wait.Done()
			for {
				port, ok := <-ports
				if !ok {
					break
				}
				if IsOpenTCP(ip, port) {
					ret = append(ret, port)
				}
			}
		}(ports, ip)
	}
	bufsnake := make(chan int)
	for i := 0; i < port_thread; i++ {
		go func(start, end int) {
			for port := start; port < end; port++ {
				ports <- strconv.Itoa(port)
			}
			bufsnake <- 1
		}(i*(65535/port_thread), (i+1)*(65535/port_thread))
	}
	for i := 0; i < port_thread; i++ {
		<-bufsnake
	}
	close(ports)
	wait.Wait()
	return ret
}

func MPortScan(ip string) string {
	ports := ""
	m := masscan.New()
	m.SetPorts("0-65535")
	m.SetArgs(ip)
	m.SetRate("1000")
	err := m.Run()
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
		return ""
	}
	results, err := m.Parse()
	if err != nil {
		return ""
	}
	for _, result := range results {
		for _, ip := range result.Ports {
			ports += ip.Portid + ","
		}
	}
	return ports
}

func IsOpenTCP(IpAddr, Port string) bool {
	conn, err := net.DialTimeout("tcp", IpAddr+":"+Port, time.Second*1/10)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func GetWeb(wait *sync.WaitGroup, ipchan chan string) {
	defer wait.Done()
	for link := range ipchan {
		http, https := "", ""
		webwait := sync.WaitGroup{}
		webwait.Add(2)
		go func() {
			http = FastHTTP(link)
			webwait.Done()
		}()
		go func() {
			https = FastHTTPS(link)
			webwait.Done()
		}()
		webwait.Wait()
		if http != "" && len(http) > 3 {
			if http[:3] != "400" && http[:3] != "503" {
				color.Green(http)
			}
		}
		if https != "" && len(https) > 3 {
			if https[:3] != "400" && https[:3] != "503" {
				color.Green(https)
			}
		}
	}
}

func FastHTTP(url string) string {
	client := http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", "http://"+url, nil)
	req.Header.Add("Connection", "close")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		return ""
	}
	if strings.Split(url, ":")[1] == "80" {
		return strconv.Itoa(response.StatusCode) + " " + "http://" + strings.Split(url, ":")[0]
	} else {
		return strconv.Itoa(response.StatusCode) + " " + "http://" + url
	}
}

func FastHTTPS(url string) string {
	client := http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, _ := http.NewRequest("GET", "https://"+url, nil)
	req.Header.Add("Connection", "close")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.102 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		return ""
	}
	if strings.Split(url, ":")[1] == "443" {
		return strconv.Itoa(response.StatusCode) + " " + "https://" + strings.Split(url, ":")[0]
	} else {
		return strconv.Itoa(response.StatusCode) + " " + "https://" + url
	}
}

func Exist(source []string, data string) bool {
	for i := 0; i < len(source); i++ {
		if source[i] == data {
			return true
		}
	}
	return false
}
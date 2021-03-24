package config

import (
	"github.com/bufsnake/Monkey/internal/other"
	"sync"
)

type Config struct {
	File          string // 从文件获取扫描IP
	IP            string // 扫描单一IP
	Thread        int    // 扫描线程
	Blast         bool   // 是否进行爆破
	Timeout       int    // http超时设置
	Output        string // 输出文件名称 默认为当前时间(2020_03_04_22_33_14.csv)
	Masscan       bool   // 使用Masscan扫描
	SocketThreads int    // socket 链接线程
	NoPing        bool   // 禁止ping扫描

	NmapsV      string // nmap -sV --version-intensity 0~9
	Port        string // nmap -p 0-65535
	MasscanRate string // masscan --rate 1000
}

var Urldatas [][]string
var urldata_l sync.Mutex

func AppendUrls(data []string) {
	urldata_l.Lock()
	defer urldata_l.Unlock()
	Urldatas = append(Urldatas, data)
}

var Blastdatas []string
var blastdatas_l sync.Mutex

func AppendBlast(data string) {
	blastdatas_l.Lock()
	defer blastdatas_l.Unlock()
	Blastdatas = append(Blastdatas, data)
}

var AllPort []string
var AllPort_l sync.Mutex

func AppendAllPort(data string) {
	AllPort_l.Lock()
	defer AllPort_l.Unlock()
	if !other.Exist(AllPort, data) {
		AllPort = append(AllPort, data)
	}
}

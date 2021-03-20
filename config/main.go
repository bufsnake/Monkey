package config

import "sync"

type Config struct {
	File    string // 从文件获取扫描IP
	IP      string // 扫描单一IP
	Thread  int    // 扫描线程
	Timeout int    // http超时设置

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

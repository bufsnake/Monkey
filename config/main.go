package config

type Config struct {
	Target      string // 单个目标
	Targets     string // 多个目标
	Thread      int    // 扫描线程
	Blast       bool   // 是否进行爆破
	Timeout     int    // http超时设置
	NmapSV      string // nmap详细程度，默认为2，数值越大，耗时越久 nmap -sV --version-intensity 0~9
	Port        string // 指定端口 nmap -p 0-65535
	MasScanRate string // masscan 扫描速率 masscan --rate 1000
}

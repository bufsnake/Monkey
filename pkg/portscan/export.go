package portscan

import (
	"github.com/bufsnake/Monkey/config"
)

type Service struct {
	IP        string
	Port      string
	Protocol  string
	Version   string
	ServiceFP string // 终端不输出，输出到CSV中
}

type portscan interface {
	PortScan(ip string) error
	GetService() []Service
}

func NewPortScan(config config.Config) portscan {
	return &mas_scan{conf: config, port: config.Port, nmapsv: config.NmapSV, masscan_rate: config.MasScanRate}
}

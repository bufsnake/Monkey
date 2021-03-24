package findport

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

func NewPortScanX(config config.Config) portscan {
	if config.Masscan {
		return &mas_scan{conf: config, port: config.Port, nmapsv: config.NmapsV, masscan_rate: config.MasscanRate}
	}
	return &socket{conf: config, port: config.Port, nmapsv: config.NmapsV, masscan_rate: config.MasscanRate}
}

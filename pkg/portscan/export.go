package findport

import "github.com/bufsnake/Monkey/config"

type findport struct {
	port         string
	nmapsv       string
	masscan_rate string
	masscan_port string
	Services     []Service
}

type Service struct {
	IP        string
	Port      string
	Protocol  string
	Version   string
	ServiceFP string
}

func NewPortScan(config config.Config) (*findport, error) {
	return &findport{port: config.Port, nmapsv: config.NmapsV, masscan_rate: config.MasscanRate}, nil
}

func (f *findport) PortScan(ip string) error {
	err := f.getport(ip)
	if err != nil {
		return err
	}
	err = f.getsevice(ip)
	if err != nil {
		return err
	}
	return nil
}

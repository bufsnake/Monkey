package blasting

import (
	"strconv"
)

type Blast interface {
	Connect() (string, error)
	Info() string
}

func NewBlast(service, ip, port string) (Blast, error) {
	atoi, err := strconv.Atoi(port)
	if err != nil {
		return nil, err
	}
	if atoi < 0 || atoi > 65535 {
		return nil, err
	}
	switch service {
	case "ssh", "ssl/ssh":
		return &clissh{ip: ip, port: port}, nil
	case "ftp", "ssl/ftp":
		return &cliftp{ip: ip, port: port}, nil
	case "mysql", "nagios-nsca", "ssl/mysql", "ssl/nagios-nsca":
		return &climysql{ip: ip, port: port}, nil
	case "ms-sql-s", "ssl/ms-sql-s":
		return &mssql{ip: ip, port: port}, nil
	case "redis", "ssl/redis":
		return &cliredis{ip: ip, port: port}, nil
	case "mongodb", "ssl/mongodb":
		return &mongodb{ip: ip, port: port}, nil
	case "postgresql", "ssl/postgresql":
		return &postgresql{ip: ip, port: port}, nil
	case "microsoft-ds", "ssl/microsoft-ds":
		return &clismb{ip: ip, port: port}, nil
	case "snmp", "ssl/snmp":
		return &snmp{ip: ip, port: port}, nil
	case "docker", "ssl/docker":
		return &docker{ip: ip, port: port}, nil
	case "zookeeper", "ssl/zookeeper":
		return &zookeeper{ip: ip, port: port}, nil
	case "memcached", "ssl/memcached":
		return &memcached{ip: ip, port: port}, nil
	default:
		return nil, nil
	}
}

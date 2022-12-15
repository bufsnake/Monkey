package blasting

import (
	"errors"
	"strconv"
	"strings"
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
	switch strings.ToUpper(service) {
	case "SSH", "SSL/SSH":
		return &clissh{ip: ip, port: port}, nil
	case "FTP", "SSL/FTP":
		return &cliftp{ip: ip, port: port}, nil
	case "MYSQL", "NAGIOS-NSCA", "SSL/MYSQL", "SSL/NAGIOS-NSCA":
		return &climysql{ip: ip, port: port}, nil
	case "MS-SQL-S", "SSL/MS-SQL-S":
		return &mssql{ip: ip, port: port}, nil
	case "REDIS", "SSL/REDIS":
		return &redis_cli{ip: ip, port: port}, nil
	case "MONGOD", "SSL/MONGOD", "MONGODB", "SSL/MONGODB":
		return &mongodb{ip: ip, port: port}, nil
	case "POSTGRESQL", "SSL/POSTGRESQL":
		return &postgresql{ip: ip, port: port}, nil
	case "MICROSOFT-DS", "SSL/MICROSOFT-DS":
		return &clismb{ip: ip, port: port}, nil
	case "SNMP", "SSL/SNMP":
		return &snmp{ip: ip, port: port}, nil
	case "DOCKER", "SSL/DOCKER":
		return &docker{ip: ip, port: port}, nil
	case "ZOOKEEPER", "SSL/ZOOKEEPER":
		return &zookeeper{ip: ip, port: port}, nil
	case "MEMCACHE", "SSL/MEMCACHE":
		return &memcached{ip: ip, port: port}, nil
	case "MS-WBT-SERVER", "SSL/MS-WBT-SERVER":
		return &rdpcli{ip: ip, port: port}, nil
	case "TELNET", "SSL/TELNET":
		return &telnetcli{ip: ip, port: port}, nil
	case "ORACLE", "SSL/ORACLE":
		return &oraclecli{ip: ip, port: port}, nil
	default:
		return nil, errors.New("found no protocol")
	}
}

package blasting

import (
	"github.com/gosnmp/gosnmp"
	"strconv"
	"time"
)

type snmp struct {
	ip   string
	port string
}

func (s *snmp) Info() string {
	return "weak"
}

func (s *snmp) Connect() (string, error) {
	gosnmp.Default.Target = s.ip
	temp, err := strconv.Atoi(s.port)
	gosnmp.Default.Port = uint16(temp)
	gosnmp.Default.Community = "public"
	gosnmp.Default.Timeout = 10 * time.Second
	err = gosnmp.Default.Connect()
	if err != nil {
		return "", err
	}
	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
	_, err = gosnmp.Default.Get(oids)
	if err != nil {
		return "", err
	}
	return "snmp://public@" + s.ip + ":" + s.port, nil
}

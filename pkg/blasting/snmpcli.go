package blasting

import (
	"errors"
	"github.com/gosnmp/gosnmp"
	"strconv"
	"time"
)

type snmp struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *snmp) Info() string {
	return "weak"
}

func (s *snmp) Connect() (string, error) {
	gosnmp.Default.Target = s.ip
	temp, err := strconv.Atoi(s.port)
	gosnmp.Default.Port = uint16(temp)
	gosnmp.Default.Community = "public"
	gosnmp.Default.Timeout = 2 * time.Second

	err = gosnmp.Default.Connect()
	if err == nil {
		oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
		_, err := gosnmp.Default.Get(oids)
		if err == nil {
			return "snmp://public@" + s.ip + ":" + s.port, nil
		}
	}
	return "", errors.New("snmp weak password test finish,but no password found")
}

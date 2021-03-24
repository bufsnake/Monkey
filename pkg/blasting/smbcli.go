package blasting

import (
	"errors"
	"github.com/stacktitan/smb/smb"
	"strconv"
	"strings"
	"time"
)

type clismb struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *clismb) Info() string {
	return "weak"
}
func (s *clismb) Connect() (string, error) {
	userdict := []string{"administrator", "admin", "test", "user", "manager", "webadmin", "guest", "db2admin"}
	passdict := []string{"%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123!@#", "admin", "admin123", "admin@123", "admin#123", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "postgres", "q1w2e3r4", "redhat", "user", "mysql", "apache", ""}
	dict := make([]string, 0)
	for i := 0; i < len(userdict); i++ {
		for j := 0; j < len(passdict); j++ {
			if !inintslice(dict, userdict[i]+":bufsnake:"+strings.Replace(passdict[j], "%user%", userdict[i], -1)) {
				dict = append(dict, userdict[i]+":bufsnake:"+strings.Replace(passdict[j], "%user%", userdict[i], -1))
			}
		}
	}
	dictchan := make(chan string, 10)
	dictlen := len(dict)
	fin := make(chan string)
	for i := 0; i < len(dict); i += 50 {
		if i+50 > len(dict) {
			go smbdict(dict[i:dictlen], dictchan)
			break
		}
		go smbdict(dict[i:i+50], dictchan)
	}
	for i := 0; i < 20; i++ {
		go smbconnect(dictchan, fin, s.ip, s.port)
		<-time.After(1 * time.Second / 1000)
	}
	for i := 0; i < dictlen; i++ {
		temp := <-fin
		if temp != "" {
			return "smb://" + temp + "@" + s.ip + ":" + s.port, nil
		}
	}
	return "", errors.New("smb weak password test finish,but no password found")
}

func smbdict(dict []string, dictchan chan string) {
	for i := 0; i < len(dict); i++ {
		dictchan <- dict[i]
	}
}

func smbconnect(dictchan, fin chan string, host string, port string) {
	for dict := range dictchan {
		user := strings.Split(dict, ":bufsnake:")[0]
		password := strings.Split(dict, ":bufsnake:")[1]
		temp, _ := strconv.Atoi(port)
		options := smb.Options{
			Host:        host,
			Port:        temp,
			User:        user,
			Password:    password,
			Domain:      "",
			Workstation: "",
		}
		session, err := smb.NewSession(options, false)
		if err == nil {
			session.Close()
			if session.IsAuthenticated {
				fin <- user + ":" + password
			}
		}
		fin <- ""
	}
}

package blasting

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"strings"
	"time"
)

type mongodb struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *mongodb) Info() string {
	return "weak"
}
func (m *mongodb) Connect() (string, error) {
	session, err := mgo.DialWithTimeout(fmt.Sprintf("mongodb://%s:%s/%s", m.ip, m.port, "test"), 10*time.Second)
	if err == nil {
		_, err := session.DatabaseNames()
		if err == nil {
			session.Close()
			return "mongodb://未授权@" + m.ip + ":" + m.port, nil
		}
		session.Close()
	}
	userdict := []string{"admin", "test", "system", "web"}
	passdict := []string{"admin", "mongodb", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "Passw0rd", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "admin123", "admin888", "administrator", "administrator123", "mongodb123", "mongodbpass", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "user", "web", ""}
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
			go mongodbdict(dict[i:dictlen], dictchan)
			break
		}
		go mongodbdict(dict[i:i+50], dictchan)
	}
	for i := 0; i < 5; i++ {
		go mongodbconnect(dictchan, fin, m.ip, m.port)
		<-time.After(1 * time.Second / 1000)
	}
	for i := 0; i < dictlen; i++ {
		temp := <-fin
		if temp != "" {
			return "mongodb://" + temp + "@" + m.ip + ":" + m.port, nil
		}
	}
	return "", errors.New("mongodb weak password test finish,but no password found")
}

func mongodbdict(dict []string, dictchan chan string) {
	for i := 0; i < len(dict); i++ {
		dictchan <- dict[i]
	}
}

func mongodbconnect(dictchan, fin chan string, host string, port string) {
	for dict := range dictchan {
		user := strings.Split(dict, ":bufsnake:")[0]
		password := strings.Split(dict, ":bufsnake:")[1]
		session, err := mgo.DialWithTimeout(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", user, password, host, port, "test"), 10*time.Second)
		if err == nil {
			defer session.Close()
			err = session.Ping()
			if err == nil {
				fin <- user + ":" + password
			}
		}
		fin <- ""
	}
}

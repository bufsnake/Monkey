package blasting

import (
	"errors"
	"github.com/ziutek/mymysql/mysql"
	"strings"
	"time"

	_ "github.com/ziutek/mymysql/native"
)

type climysql struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *climysql) Info() string {
	return "weak"
}
func (m *climysql) Connect() (string, error) {
	userdict := []string{"root"}
	passdict := []string{"%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123", "admin888", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "123456qwerty", "password123", "1q2w3e", "okmnji", "test123", "test12345", "test123456", "q1w2e3r4", "mysql", "web", "%username%", "%null%", "123", "1234", "12345", "123456", "admin", "pass", "password", "!null!", "", "!user!", "", "1234567", "7654321", "abc123", "111111", "123321", "123123", "12345678", "123456789", "000000", "888888", "654321", "987654321", "147258369", "123asd", "qwer123", "P@ssw0rd", "root3306", "Q1W2E3b3", ""}
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
			go mysqldict(dict[i:dictlen], dictchan)
			break
		}
		go mysqldict(dict[i:i+50], dictchan)
	}
	for i := 0; i < 5; i++ {
		go mysqlconnect(dictchan, fin, m.ip, m.port)
		<-time.After(1 * time.Second / 1000)
	}
	for i := 0; i < dictlen; i++ {
		temp := <-fin
		if temp != "" {
			return "mysql://" + temp + "@" + m.ip + ":" + m.port, nil
		}
	}
	return "", errors.New("mysql weak password test finish,but no password found")
}

func mysqldict(dict []string, dictchan chan string) {
	for i := 0; i < len(dict); i++ {
		dictchan <- dict[i]
	}
}

func mysqlconnect(dictchan, fin chan string, host string, port string) {
	for dict := range dictchan {
		user := strings.Split(dict, ":bufsnake:")[0]
		password := strings.Split(dict, ":bufsnake:")[1]
		conn := mysql.New("tcp", "", host+":"+port, user, password, "mysql")
		conn.SetTimeout(10 * time.Second)
		err := conn.Connect()
		if err != nil {
			conn.Close()
			fin <- ""
		} else {
			conn.Close()
			fin <- user + ":" + password
		}

	}
}

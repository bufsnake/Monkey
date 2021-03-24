package blasting

import (
	"errors"
	"github.com/jlaffaye/ftp"
	"strings"
	"time"
)

type cliftp struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *cliftp) Info() string {
	return "weak"
}

func (f *cliftp) Connect() (string, error) {
	userdict := []string{"anonymous", "administrator", "ftp", "test", "admin", "web"}
	passdict := []string{"%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123", "admin888", "administrator", "administrator123", "ftp", "ftppass", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "user", "mysql", "web", ""}
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
			go ftpdict(dict[i:dictlen], dictchan)
			break
		}
		go ftpdict(dict[i:i+50], dictchan)
	}
	for i := 0; i < 20; i++ {
		go ftpconnect(dictchan, fin, f.ip, f.port)
		<-time.After(1 * time.Second / 1000)
	}
	for i := 0; i < dictlen; i++ {
		temp := <-fin
		if temp != "" {
			return "ftp://" + temp + "@" + f.ip + ":" + f.port, nil
		}
	}
	return "", errors.New("ftp weak password test finish,but no password found")
}

func ftpdict(dict []string, dictchan chan string) {
	for i := 0; i < len(dict); i++ {
		dictchan <- dict[i]
	}
}

func ftpconnect(dictchan, fin chan string, host string, port string) {
	for dict := range dictchan {
		user := strings.Split(dict, ":bufsnake:")[0]
		password := strings.Split(dict, ":bufsnake:")[1]
		c, err := ftp.Dial(host+":"+port, ftp.DialWithTimeout(10*time.Second))
		if err == nil {
			err = c.Login(user, password)
			if err == nil {
				if err := c.Quit(); err == nil {
					fin <- user + ":" + password
				} else {
					fin <- ""
				}
			} else {
				fin <- ""
			}
		} else {
			fin <- ""
		}
	}
}

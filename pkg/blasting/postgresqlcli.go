package blasting

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type postgresql struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *postgresql) Info() string {
	return "weak"
}
func (p *postgresql) Connect() (string, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s/%s?sslmode=%s", p.ip, p.port, "postgres", "disable"))
	if err == nil {
		_, err := db.Query("fjdaslfjaskldfjasl")
		if !strings.Contains(err.Error(), "authentication") && strings.Contains(err.Error(), "syntax error at or near") {
			err = db.Close()
			return "postgresql://未授权@" + p.ip + ":" + p.port, nil
		}
		err = db.Close()
	}
	userdict := []string{"postgres", "test", "admin", "web"}
	passdict := []string{"admin", "Passw0rd", "postgres", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "admin123", "admin888", "administrator", "administrator123", "root123", "ftp", "ftppass", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "user", "web", ""}
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
			go postgresqldict(dict[i:dictlen], dictchan)
			break
		}
		go postgresqldict(dict[i:i+50], dictchan)
	}
	for i := 0; i < 5; i++ {
		go postgresqlconnect(dictchan, fin, p.ip, p.port)
		<-time.After(1 * time.Second / 1000)
	}
	for i := 0; i < dictlen; i++ {
		temp := <-fin
		if temp != "" {
			return "postgresql://" + temp + "@" + p.ip + ":" + p.port, nil
		}
	}
	return "", errors.New("postgresql weak password test finish,but no password found")
}

func postgresqldict(dict []string, dictchan chan string) {
	for i := 0; i < len(dict); i++ {
		dictchan <- dict[i]
	}
}

func postgresqlconnect(dictchan, fin chan string, host string, port string) {
	for dict := range dictchan {
		user := strings.Split(dict, ":bufsnake:")[0]
		password := strings.Split(dict, ":bufsnake:")[1]
		dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, "postgres", "disable")
		db, err := sql.Open("postgres", dataSourceName)
		if err == nil {
			defer db.Close()
			err = db.Ping()
			if err == nil {
				fin <- user + ":" + password
			}
		}
		fin <- ""
	}
}

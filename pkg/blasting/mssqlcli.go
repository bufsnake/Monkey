package blasting

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type mssql struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *mssql) Info() string {
	return "weak"
}
func (m *mssql) Connect() (string, error) {
	userdict := []string{"sa"}
	passdict := []string{"admin", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "qweasdzxc", "admin123", "admin888", "administrator", "administrator123", "sa123", "ftp", "ftppass", "123456", "password", "12345", "1234", "sa", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "sqlpass", "sql123", "sqlserver", "web", ""}
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
			go mssqldict(dict[i:dictlen], dictchan)
			break
		}
		go mssqldict(dict[i:i+50], dictchan)
	}
	for i := 0; i < 5; i++ {
		go mssqlconnect(dictchan, fin, m.ip, m.port)
		<-time.After(1 * time.Second / 1000)
	}
	for i := 0; i < dictlen; i++ {
		temp := <-fin
		if temp != "" {
			return "mssql://" + temp + "@" + m.ip + ":" + m.port, nil
		}
	}
	return "", errors.New("mssql weak password test finish,but no password found")
}

func mssqldict(dict []string, dictchan chan string) {
	for i := 0; i < len(dict); i++ {
		dictchan <- dict[i]
	}
}

func mssqlconnect(dictchan, fin chan string, host string, port string) {
	for dict := range dictchan {
		user := strings.Split(dict, ":bufsnake:")[0]
		password := strings.Split(dict, ":bufsnake:")[1]
		dataSourceName := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s", host, port, user, password, "master")
		db, err := sql.Open("mssql", dataSourceName)
		if err != nil {
			fin <- ""
		} else {
			defer db.Close()
			err := db.Ping()
			if err == nil {
				fin <- user + ":" + password
			} else {
				fin <- ""
			}
		}
	}
}

package blasting

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
	_ "github.com/ziutek/mymysql/native"
	"log"
	"strings"
	"sync"
	"time"
)

type oraclecli struct {
	ip   string
	port string
}

func (m *oraclecli) Info() string {
	return "weak"
}

func (m *oraclecli) Connect() (string, error) {
	usernames := []string{"sys", "system", "admin", "test", "web", "orcl"}
	passwords := []string{"R0ck9", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin", "sys", "system", "oracle", "dbadmin", "qweasdzxc", "admin123", "admin888", "administrator", "administrator123", "root123", "123456", "password", "12345", "root", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4"}
	wait := sync.WaitGroup{}
	messages := make(chan message, 600)
	fin := make(chan message)
	for i := 0; i < 5; i++ {
		wait.Add(1)
		go m.check(&wait, messages, fin)
	}
	for _, u := range usernames {
		for _, p := range passwords {
			messages <- message{
				user: u,
				pass: p,
			}
		}
	}
	go func() {
		close(messages)
		wait.Wait()
		fin <- message{user: "error"}
	}()
	select {
	case <-time.After(5 * time.Minute):
		return "", errors.New("oracle weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("oracle weak password test finish,but no password found")
		}
		return "oracle://" + mess.user + ":" + mess.pass + "@" + m.ip + ":" + m.port, nil
	}
}

func (m *oraclecli) check(wg *sync.WaitGroup, messages, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		var db *sql.DB
		var err error
		for _, service := range []string{"ORCL", "XE", "ORACLE"} {
			db, err = sql.Open("oracle", fmt.Sprintf("oracle://%s:%s@%s:%s/%s", message_.user, message_.pass, m.ip, m.port, service))
			if err == nil {
				break
			}
		}
		db.SetConnMaxLifetime(5 * time.Second)
		db.SetMaxIdleConns(0)
		err = db.Ping()
		if err != nil {
			log.Println(err)
			db.Close()
			continue
		}
		db.Close()
		fin <- message_
	}
}

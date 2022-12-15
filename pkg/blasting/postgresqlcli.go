package blasting

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
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

	usernames := []string{"postgres", "test", "admin", "web"}
	passwords := []string{"admin", "Passw0rd", "postgres", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "admin123", "admin888", "administrator", "administrator123", "root123", "ftp", "ftppass", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "user", "web", ""}
	wait := sync.WaitGroup{}
	messages := make(chan message, 600)
	fin := make(chan message)
	for i := 0; i < 5; i++ {
		wait.Add(1)
		go p.check(&wait, messages, fin)
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
		return "", errors.New("postgresql weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("postgresql weak password test finish,but no password found")
		}
		return "postgresql://" + mess.user + ":" + mess.pass + "@" + p.ip + ":" + p.port, nil
	}
}

func (p *postgresql) check(wg *sync.WaitGroup, messages, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", message_.user, message_.pass, p.ip, p.port, "postgres", "disable"))
		if err != nil {
			log.Println(err)
			continue
		}
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

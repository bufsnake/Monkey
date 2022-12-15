package blasting

import (
	"errors"
	"github.com/ziutek/mymysql/mysql"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/ziutek/mymysql/native"
)

type climysql struct {
	ip   string
	port string
}

func (m *climysql) Info() string {
	return "weak"
}

func (m *climysql) Connect() (string, error) {
	usernames := []string{"root"}
	passwords := []string{"%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123", "admin888", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "123456qwerty", "password123", "1q2w3e", "okmnji", "test123", "test12345", "test123456", "q1w2e3r4", "mysql", "web", "%username%", "%null%", "123", "1234", "12345", "123456", "admin", "pass", "password", "!null!", "", "!user!", "", "1234567", "7654321", "abc123", "111111", "123321", "123123", "12345678", "123456789", "000000", "888888", "654321", "987654321", "147258369", "123asd", "qwer123", "P@ssw0rd", "root3306", "Q1W2E3b3", ""}
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
		return "", errors.New("mysql weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("mysql weak password test finish,but no password found")
		}
		return "mysql://" + mess.user + ":" + mess.pass + "@" + m.ip + ":" + m.port, nil
	}
}

func (m *climysql) check(wg *sync.WaitGroup, messages, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		conn := mysql.New("tcp", "", m.ip+":"+m.port, message_.user, message_.pass, "mysql")
		conn.SetTimeout(5 * time.Second)
		err := conn.Connect()
		if err != nil {
			log.Println(err)
			continue
		}
		conn.Close()
		fin <- message_
	}
}

package blasting

import (
	"errors"
	"github.com/bufsnake/Monkey/pkg/grdp"
	"log"
	"strings"
	"sync"
	"time"
)

type rdpcli struct {
	ip   string
	port string
}

func (s *rdpcli) Info() string {
	return "weak"
}

func (m *rdpcli) Connect() (string, error) {
	usernames := []string{"woziniyaqi", "administrator", "admin", "test", "user", "manager", "webadmin", "guest", "db2admin"}
	passwords := []string{"R0ck9", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123!@#", "admin", "admin123", "admin@123", "admin#123", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "postgres", "q1w2e3r4", "redhat", "user", "mysql", "apache"}
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
		return "", errors.New("rdp weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("rdp weak password test finish,but no password found")
		}
		return "rdp://" + mess.user + ":" + mess.pass + "@" + m.ip + ":" + m.port, nil
	}
}

func (m *rdpcli) check(wg *sync.WaitGroup, messages chan message, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		err := grdp.Login(m.ip+":"+m.port, "", message_.user, message_.pass)
		if err != nil {
			log.Println("rdp cli blast error", message_.user, message_.pass, err)
			continue
		}
		fin <- message_
	}
}

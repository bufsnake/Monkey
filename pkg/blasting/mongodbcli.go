package blasting

import (
	"errors"
	"fmt"
	"gopkg.in/mgo.v2"
	"log"
	"strings"
	"sync"
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
	session, err := mgo.DialWithTimeout(fmt.Sprintf("mongodb://%s:%s/%s", m.ip, m.port, "test"), 2*time.Second)
	if err == nil {
		_, err := session.DatabaseNames()
		if err == nil {
			session.Close()
			return "mongodb://未授权@" + m.ip + ":" + m.port, nil
		}
		session.Close()
	}

	usernames := []string{"admin", "test", "system", "web"}
	passwords := []string{"123456", "admin", "mongodb", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "Passw0rd", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "admin123", "admin888", "administrator", "administrator123", "mongodb123", "mongodbpass", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "user", "web", ""}
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
		return "", errors.New("mongodb weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("mongodb weak password test finish,but no password found")
		}
		return "mongodb://" + mess.user + ":" + mess.pass + "@" + m.ip + ":" + m.port, nil
	}
}

func (m *mongodb) check(wg *sync.WaitGroup, messages chan message, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		session, err := mgo.DialWithTimeout(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", message_.user, message_.pass, m.ip, m.port, "admin"), 10*time.Second)
		if err != nil {
			log.Println(err)
			continue
		}
		err = session.Ping()
		if err != nil {
			log.Println(err)
			session.Close()
			continue
		}
		session.Close()
		fin <- message_
	}
}

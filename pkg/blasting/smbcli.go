package blasting

import (
	"errors"
	"github.com/stacktitan/smb/smb"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type clismb struct {
	ip   string
	port string
}

func (s *clismb) Info() string {
	return "weak"
}

func (s *clismb) Connect() (string, error) {
	usernames := []string{"administrator", "admin", "test", "user", "manager", "webadmin", "guest", "db2admin"}
	passwords := []string{"%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123!@#", "admin", "admin123", "admin@123", "admin#123", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "postgres", "q1w2e3r4", "redhat", "user", "mysql", "apache", ""}
	wait := sync.WaitGroup{}
	messages := make(chan message, 600)
	fin := make(chan message)
	for i := 0; i < 5; i++ {
		wait.Add(1)
		go s.check(&wait, messages, fin)
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
		return "", errors.New("smb weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("smb weak password test finish,but no password found")
		}
		return "smb://" + mess.user + ":" + mess.pass + "@" + s.ip + ":" + s.port, nil
	}
}

func (s *clismb) check(wg *sync.WaitGroup, messages, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		port, err := strconv.Atoi(s.port)
		if err != nil {
			log.Println(s.ip, s.port, err)
			continue
		}
		options := smb.Options{
			Host:        s.ip,
			Port:        port,
			User:        message_.user,
			Password:    message_.pass,
			Domain:      "",
			Workstation: "",
		}
		session, err := smb.NewSession(options, false)
		if err != nil {
			log.Println(err)
			continue
		}
		session.Close()
		if session.IsAuthenticated {
			fin <- message_
		}
	}
}

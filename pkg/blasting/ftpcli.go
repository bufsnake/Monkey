package blasting

import (
	"errors"
	"github.com/jlaffaye/ftp"
	"log"
	"strings"
	"sync"
	"time"
)

type cliftp struct {
	ip   string
	port string
}

func (f *cliftp) Info() string {
	return "weak"
}

func (f *cliftp) Connect() (string, error) {
	usernames := []string{"anonymous", "administrator", "ftp", "test", "admin", "web"}
	passwords := []string{"%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123", "admin888", "administrator", "administrator123", "ftp", "ftppass", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "user", "mysql", "web", ""}
	wait := sync.WaitGroup{}
	messages := make(chan message, 600)
	fin := make(chan message)
	for i := 0; i < 5; i++ {
		wait.Add(1)
		go f.check(&wait, messages, fin)
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
		return "", errors.New("ftp weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("ftp weak password test finish,but no password found")
		}
		return "ftp://" + mess.user + ":" + mess.pass + "@" + f.ip + ":" + f.port, nil
	}
}

func (f *cliftp) check(wg *sync.WaitGroup, messages, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		c, err := ftp.Dial(f.ip+":"+f.port, ftp.DialWithTimeout(5*time.Second))
		if err != nil {
			log.Println(err)
			continue
		}
		err = c.Login(message_.user, message_.pass)
		if err != nil {
			log.Println(err)
			continue
		}
		err = c.Quit()
		if err != nil {
			log.Println(err)
			continue
		}
		fin <- message_
	}
}

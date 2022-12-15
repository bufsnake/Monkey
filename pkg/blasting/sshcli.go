package blasting

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type clissh struct {
	ip   string
	port string
}

func (s *clissh) Info() string {
	return "weak"
}

func (s *clissh) Connect() (string, error) {
	usernames := []string{"root", "jyhcg", "d5000", "ns5000", "test", "oracle", "admin", "user", "postgres", "mysql", "backup", "guest", "system", "web", "tomcat", "michael", "upload", "alex", "sys", "sales", "linux", "ftp", "temp", "nagios", "user1", "www", "test1", "eSER!@#"}
	passwords := []string{"%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123!@#", "admin", "admin123", "admin@123", "admin#123", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "postgres", "q1w2e3r4", "redhat", "user", "mysql", "apache", "d5000", "ns5000", "jyhcg", ""}
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
		return "", errors.New("ssh weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("ssh weak password test finish,but no password found")
		}
		return "ssh://" + mess.user + ":" + mess.pass + "@" + s.ip + ":" + s.port, nil
	}
}

func (s *clissh) check(wg *sync.WaitGroup, messages, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		var (
			addr         string
			clientConfig *ssh.ClientConfig
			err          error
		)
		clientConfig = &ssh.ClientConfig{
			User:    message_.user,
			Auth:    []ssh.AuthMethod{ssh.Password(message_.pass)},
			Timeout: 5 * time.Second,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
		clientConfig.Ciphers = []string{"aes128-cbc", "chacha20-poly1305@openssh.com", "aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "aes256-gcm@openssh.com"}
		addr = fmt.Sprintf("%s:%s", s.ip, s.port)
		_, err = ssh.Dial("tcp", addr, clientConfig)
		if err != nil {
			log.Println(err)
			continue
		}
		fin <- message_
	}
}

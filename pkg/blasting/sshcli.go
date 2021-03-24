package blasting

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"strings"
	"time"
)

type clissh struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *clissh) Info() string {
	return "weak"
}

func (s *clissh) Connect() (string, error) {
	userdict := []string{"root", "jyhcg", "d5000", "ns5000", "test", "oracle", "admin", "user", "postgres", "mysql", "backup", "guest", "system", "web", "tomcat", "michael", "upload", "alex", "sys", "sales", "linux", "ftp", "temp", "nagios", "user1", "www", "test1", "eSER!@#"}
	passdict := []string{"R0ck9", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "admin123!@#", "admin", "admin123", "admin@123", "admin#123", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "postgres", "q1w2e3r4", "redhat", "user", "mysql", "apache", "d5000", "ns5000", "jyhcg", ""}
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
			go sshdict(dict[i:dictlen], dictchan)
			break
		}
		go sshdict(dict[i:i+50], dictchan)
	}
	for i := 0; i < 20; i++ {
		go sshconnect(dictchan, fin, s.ip, s.port)
		<-time.After(1 * time.Second / 1000)
	}
	for i := 0; i < dictlen; i++ {
		temp := <-fin
		if temp != "" {
			return "ssh://" + temp + "@" + s.ip + ":" + s.port, nil
		}
	}
	return "", errors.New("ssh weak password test finish,but no password found")
}

func sshdict(dict []string, dictchan chan string) {
	for i := 0; i < len(dict); i++ {
		dictchan <- dict[i]
	}
}

func sshconnect(dictchan, fin chan string, host string, port string) {
	for dict := range dictchan {
		user := strings.Split(dict, ":bufsnake:")[0]
		password := strings.Split(dict, ":bufsnake:")[1]
		var (
			auth         []ssh.AuthMethod
			addr         string
			clientConfig *ssh.ClientConfig
			err          error
		)
		auth = make([]ssh.AuthMethod, 0)
		auth = append(auth, ssh.Password(password))
		clientConfig = &ssh.ClientConfig{
			User:    user,
			Auth:    auth,
			Timeout: 10 * time.Second,
			HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			},
		}
		clientConfig.Ciphers = []string{"aes128-cbc", "aes192-ctr", "aes256-ctr"}
		addr = fmt.Sprintf("%s:%s", host, port)
		_, err = ssh.Dial("tcp", addr, clientConfig)
		if err != nil {
			if strings.Contains(err.Error(), "no common algorithm for client to server cipher") {
				fmt.Println("\n" + err.Error())
			}
		}
		if err == nil {
			fin <- user + ":" + password
		}
		fin <- ""
	}
}

func inintslice(haystack []string, needle string) bool {
	for _, e := range haystack {
		if e == needle {
			return true
		}
	}

	return false
}

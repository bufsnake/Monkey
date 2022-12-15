package blasting

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type telnetcli struct {
	ip   string
	port string
}

func (s *telnetcli) Info() string {
	return "weak"
}

func (m *telnetcli) Connect() (string, error) {
	usernames := []string{"admin", "root", "guest", "ftp", "www"}
	passwords := []string{"R0ck9", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "qweasdzxc", "Passw0rd", "admin", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "apache", "qwer1234"}
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
		return "", errors.New("telnet weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("telnet weak password test finish,but no password found")
		}
		return "telnet://" + mess.user + ":" + mess.pass + "@" + m.ip + ":" + m.port, nil
	}
}

func (m *telnetcli) check(wg *sync.WaitGroup, messages chan message, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		t := telnet{
			IP:               m.ip,
			Port:             m.port,
			IsAuthentication: true,
			UserName:         message_.user,
			Password:         message_.pass,
		}
		telnet_, err := t.Telnet(10)
		if err != nil {
			log.Println("rdp cli blast error", message_.user, message_.pass, err)
			continue
		}
		if telnet_ {
			fin <- message_
		}
	}
}

type telnet struct {
	IP               string
	Port             string
	IsAuthentication bool
	UserName         string
	Password         string
}

func (this *telnet) Telnet(timeout int) (bool, error) {
	raddr := this.IP + ":" + this.Port
	conn, err := net.DialTimeout("tcp", raddr, time.Duration(timeout)*time.Second)
	if nil != err {
		log.Print("pkg: model, func: Telnet, method: net.DialTimeout, errInfo:", err)
		return false, err
	}
	defer conn.Close()
	if false == this.telnetProtocolHandshake(conn) {
		//log.Print("pkg: model, func: Telnet, method: this.telnetProtocolHandshake, errInfo: telnet protocol handshake failed!!!")
		return false, err
	}
	return true, err
}

func (this *telnet) telnetProtocolHandshake(conn net.Conn) bool {
	var buf [4096]byte
	var n int
	n, err := conn.Read(buf[0:])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake1, method: conn.Read, errInfo:", err)
		return false
	}

	buf[0] = 0xff
	buf[1] = 0xfc
	buf[2] = 0x25
	buf[3] = 0xff
	buf[4] = 0xfe
	buf[5] = 0x01
	n, err = conn.Write(buf[0:6])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake2, method: conn.Write, errInfo:", err)
		return false
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake3, method: conn.Read, errInfo:", err)
		return false
	}

	buf[0] = 0xff
	buf[1] = 0xfe
	buf[2] = 0x03
	buf[3] = 0xff
	buf[4] = 0xfc
	buf[5] = 0x27
	n, err = conn.Write(buf[0:6])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake4, method: conn.Write, errInfo:", err)
		return false
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake5, method: conn.Read, errInfo:", err)
		return false
	}

	//fmt.Println((buf[0:n]))
	n, err = conn.Write([]byte(this.UserName + "\r\n"))
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake6, method: conn.Write, errInfo:", err)
		return false
	}
	time.Sleep(time.Millisecond * 500)

	n, err = conn.Read(buf[0:])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake7, method: conn.Read, errInfo:", err)
		return false
	}

	n, err = conn.Write([]byte(this.Password + "\r\n"))
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake8, method: conn.Write, errInfo:", err)
		return false
	}
	time.Sleep(time.Millisecond * 2000)
	n, err = conn.Read(buf[0:])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake9, method: conn.Read, errInfo:", err)
		return false
	}
	if strings.Contains(string(buf[0:n]), "Login Failed") {
		return false
	}

	buf[0] = 0xff
	buf[1] = 0xfc
	buf[2] = 0x18

	n, err = conn.Write(buf[0:3])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake6, method: conn.Write, errInfo:", err)
		return false
	}
	n, err = conn.Read(buf[0:])
	if nil != err {
		log.Print("pkg: model, func: telnetProtocolHandshake7, method: conn.Read, errInfo:", err)
		return false
	}
	return true
}

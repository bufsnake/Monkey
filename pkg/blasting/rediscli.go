package blasting

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type redis_cli struct {
	ip   string
	port string
}

func (r *redis_cli) Info() string {
	return "weak"
}

func (r *redis_cli) Connect() (string, error) {
	conn, err := net.DialTimeout("tcp", r.ip+":"+r.port, 3*time.Second)
	if err != nil {
		log.Println("Redis Unauth Check Error", err)
	} else {
		conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
		_, err = conn.Write([]byte("info\r\n"))
		if err != nil {
			log.Println("Redis Unauth Check Error", err)
		} else {
			conn.SetReadDeadline(time.Now().Add(3 * time.Second))
			reader := bufio.NewReader(conn)
			line, _ := reader.ReadString(byte('\n'))
			lines, _ := reader.ReadString(byte('\n'))
			liness, _ := reader.ReadString(byte('\n'))
			linesss, _ := reader.ReadString(byte('\n'))
			linessss, _ := reader.ReadString(byte('\n'))
			if strings.Contains(line+lines+liness+linesss+linessss, "redis_version") {
				return "redis://未授权@" + r.ip + ":" + r.port, nil
			}
		}
		conn.Close()
	}

	//usernames := []string{"admin", "redis", "root"}
	usernames := []string{""}
	passwords := []string{"Passw0rd", "admin", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "admin123", "admin888", "administrator", "administrator123", "root123", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "user", "web"}
	wait := sync.WaitGroup{}
	messages := make(chan message, 600)
	fin := make(chan message)
	for i := 0; i < 5; i++ {
		wait.Add(1)
		go r.check(&wait, messages, fin)
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
		return "", errors.New("redis weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("redis weak password test finish,but no password found")
		}
		return "redis://" + mess.pass + "@" + r.ip + ":" + r.port, nil
	}
}

func (r *redis_cli) check(wg *sync.WaitGroup, messages, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", r.ip, r.port))
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = c.Do("AUTH", message_.pass)
		//flag := true
		//retry:
		//	if flag {
		//	} else {
		//		_, err = c.Do("AUTH", message_.user, message_.pass)
		//	}
		if err != nil {
			//if flag && strings.Contains(err.Error(), "WRONGPASS invalid username-password pair or user is disabled.") {
			//	flag = false
			//	goto retry
			//}
			log.Println(err)
			c.Close()
			continue
		}
		c.Close()
		fin <- message_
	}
}

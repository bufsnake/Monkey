package blasting

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"net"
	"strings"
	"time"
)

type cliredis struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *cliredis) Info() string {
	return "weak"
}
func (r *cliredis) Connect() (string, error) {
	conn, err := net.DialTimeout("tcp", r.ip+":"+r.port, 10*time.Second)
	if err != nil {
		log.Println("Redis Unauth Check Error", err)
	} else {
		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, err = conn.Write([]byte("info\r\n"))
		if err != nil {
			log.Println("Redis Unauth Check Error", err)
		} else {
			conn.SetReadDeadline(time.Now().Add(10 * time.Second))
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

	userdict := []string{""}
	passdict := []string{"Passw0rd", "admin", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "admin123", "admin888", "administrator", "administrator123", "root123", "123456", "password", "12345", "1234", "root", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "user", "web"}
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
			go redisdict(dict[i:dictlen], dictchan)
			break
		}
		go redisdict(dict[i:i+50], dictchan)
	}
	for i := 0; i < 5; i++ {
		go redisconnect(dictchan, fin, r.ip, r.port)
		<-time.After(1 * time.Second / 1000)
	}
	for i := 0; i < dictlen; i++ {
		temp := <-fin
		if temp != "" {
			return "redis://" + temp + "@" + r.ip + ":" + r.port, nil
		}
	}
	return "", errors.New("redis weak password test finish,but no password found")
}

func redisdict(dict []string, dictchan chan string) {
	for i := 0; i < len(dict); i++ {
		dictchan <- dict[i]
	}
}

func redisconnect(dictchan, fin chan string, host string, port string) {
	for dict := range dictchan {
		_ = strings.Split(dict, ":bufsnake:")[0]
		password := strings.Split(dict, ":bufsnake:")[1]
		c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			fin <- ""
			continue
		}
		defer c.Close()
		_, err = c.Do("AUTH", password)
		if err != nil {
			fin <- ""
			continue
		}
		fin <- password
	}
}

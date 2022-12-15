package blasting

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"
)

type memcached struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *memcached) Info() string {
	return "unau"
}

func (m *memcached) Connect() (string, error) {
	conn, err := net.DialTimeout("tcp", m.ip+":"+m.port, 3*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write([]byte("stats\r\n"))
	if err != nil {
		return "", err
	}
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	reader := bufio.NewReader(conn)
	line, _ := reader.ReadString(byte('\n'))
	if !strings.Contains(line, "STAT") {
		return "", errors.New("memcached unauth test finish,but no connect found")
	}
	return "memcached://" + m.ip + ":" + m.port, nil
}

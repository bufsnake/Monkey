package blasting

import (
	"bufio"
	"errors"
	"net"
	"strings"
	"time"
)

type zookeeper struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *zookeeper) Info() string {
	return "unau"
}
func (z *zookeeper) Connect() (string, error) {
	conn, err := net.DialTimeout("tcp", z.ip+":"+z.port, 3*time.Second)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
	_, err = conn.Write([]byte("envi\r\n"))
	if err != nil {
		return "", err
	}
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	reader := bufio.NewReader(conn)
	line, _ := reader.ReadString(byte('\n'))
	if !strings.Contains(line, "Environment") {
		return "", errors.New("zookeeper unauth test finish,but no connect found")
	}
	return "zookeeper://" + z.ip + ":" + z.port, err
}

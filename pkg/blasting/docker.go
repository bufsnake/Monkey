package blasting

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

type docker struct {
	ip   string
	port string
}

// 方便判断是啥类型的漏洞
func (s *docker) Info() string {
	return "unau"
}

func (d *docker) Connect() (string, error) {
	get, err := http.Get("http://" + d.ip + ":" + d.port + "/version")
	if err != nil {
		return "", err
	}
	defer get.Body.Close()
	all, err := ioutil.ReadAll(get.Body)
	if err != nil {
		return "", err
	}
	if strings.Contains(string(all), "Version") && strings.Contains(string(all), "Arch") && strings.Contains(string(all), "Os") {
		return "docker://" + d.ip + ":" + d.port, nil
	}
	return "", errors.New("docker unauth test finish,but no connect found")
}

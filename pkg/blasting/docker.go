package blasting

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type docker struct {
	ip   string
	port string
}

func (s *docker) Info() string {
	return "unau"
}

func (d *docker) Connect() (string, error) {
	protocol := "docker"
	req_url := fmt.Sprintf("%s://%s:%s/version", "http", d.ip, d.port)
	res, err := d.req(req_url)
	if err != nil {
		protocol = "dockers"
		log.Println(err)
		req_url = fmt.Sprintf("%s://%s:%s/version", "https", d.ip, d.port)
		res, err = d.req(req_url)
		if err != nil {
			return "", err
		}
	}
	defer res.Body.Close()
	all, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	if strings.Contains(string(all), "Version") && strings.Contains(string(all), "Arch") && strings.Contains(string(all), "Os") {
		return protocol + "://" + d.ip + ":" + d.port, nil
	}
	return "", errors.New("docker unauth test finish,but no connect found")
}

func (d *docker) req(req_url string) (*http.Response, error) {
	cli := http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest("GET", req_url, nil)
	if err != nil {
		return nil, err
	}
	return cli.Do(req)
}

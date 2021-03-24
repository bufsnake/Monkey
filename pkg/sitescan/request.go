package web

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/bufsnake/Sea/pkg/fingerprint"
	"github.com/bufsnake/Sea/pkg/useragent"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type request struct {
	url          string
	timeout      int
	title        string
	middleware   string
	xpoweredby   string
	code         int
	length       int
	header       map[string]string
	body         string
	product      string
	product_rule string
}

func (r *request) Run() error {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
			DisableKeepAlives: true,
		},
		Timeout: time.Duration(r.timeout) * time.Second,
		// 禁止301/302跳转
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	request, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		return err
	}
	request.Header.Set("USER-AGENT", useragent.RandomUserAgent())
	request.Header.Set("Connection", "close")
	request.Header.Set("rememberMe", "xxxxxxxxxxxxxxx")
	do, err := client.Do(request)
	if err != nil {
		return err
	}
	defer do.Body.Close()
	body, err := ioutil.ReadAll(do.Body)
	if err != nil {
		return err
	}
	r.length = len(body)
	r.body = string(body)
	r.code = do.StatusCode
	r.title, err = extracttitle(r.body)
	if err != nil {
		r.title = "None"
	}
	for i := 0; i < 3; i++ {
		r.title = strings.TrimLeft(r.title, " ")
		r.title = strings.TrimLeft(r.title, "\t")
		r.title = strings.TrimRight(r.title, " ")
		r.title = strings.TrimRight(r.title, "\t")
	}
	r.header = make(map[string]string)
	for key, value := range do.Header {
		val := ""
		for i := 0; i < len(value); i++ {
			val += value[i] + " "
		}
		r.header[key] = strings.ToLower(val)
	}
	r.middleware = r.GetServer()
	r.xpoweredby = r.GetXPoweredBy()
	r.getproduct()
	return nil
}

func (r *request) GetProduct() string {
	// fingerprint
	return r.product
}

func (r *request) getproduct() {
	// fingerprint
	product := fingerprint.NewFingerprint(r.url, &r.title, &r.header, &r.body)
	run, err := product.Run()
	if err != nil {
		return
	}
	if run != "" {
		rule := product.GetRule()
		temp, _ := json.Marshal(&rule)
		r.product = run
		r.product_rule = string(temp)
	}
	return
}

func (r *request) GetProductRule() string {
	return r.product_rule
}

func (r *request) GetHeader(key string) string {
	for v, k := range r.header {
		if strings.ToLower(v) == strings.ToLower(key) {
			return k
		}
	}
	return ""
}

func (r *request) GetTitle() string {
	return r.title
}

func (r *request) GetServer() string {
	return r.GetHeader("Server")
}

func (r *request) GetXPoweredBy() string {
	return r.GetHeader("X-Powered-By")
}

func (r *request) GetCode() int {
	return r.code
}

func (r *request) GetLength() int {
	return r.length
}

func (r *request) GetBody() string {
	return r.body
}

func (r *request) GetUrl() string {
	return r.url
}

// 获取网站标题
func extracttitle(body string) (string, error) {
	title := ""
	var re = regexp.MustCompile(`(?im)<\s*title.*>(.*?)<\s*/\s*title>`)
	for _, match := range re.FindAllString(body, -1) {
		title = html.UnescapeString(trimTitleTags(match))
		break
	}
	if !validUTF8([]byte(title)) {
		reader := transform.NewReader(bytes.NewReader([]byte(title)), simplifiedchinese.GBK.NewDecoder())
		d, err := ioutil.ReadAll(reader)
		if err != nil {
			return title, err
		}
		return string(d), nil
	}
	return title, nil
}

func trimTitleTags(title string) string {
	titleBegin := strings.Index(title, ">")
	titleEnd := strings.Index(title, "</")
	return title[titleBegin+1 : titleEnd]
}

func validUTF8(buf []byte) bool {
	nBytes := 0
	for i := 0; i < len(buf); i++ {
		if nBytes == 0 {
			if (buf[i] & 0x80) != 0 { //与操作之后不为0，说明首位为1
				for (buf[i] & 0x80) != 0 {
					buf[i] <<= 1 //左移一位
					nBytes++     //记录字符共占几个字节
				}

				if nBytes < 2 || nBytes > 6 { //因为UTF8编码单字符最多不超过6个字节
					return false
				}

				nBytes-- //减掉首字节的一个计数
			}
		} else { //处理多字节字符
			if buf[i]&0xc0 != 0x80 { //判断多字节后面的字节是否是10开头
				return false
			}
			nBytes--
		}
	}
	return nBytes == 0
}

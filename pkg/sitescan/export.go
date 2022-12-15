package sitescan

func NewHttpx(host string, timeout int) *httpx {
	return &httpx{url: host, timeout: timeout}
}

func NewRequest(url string, timeout int) *request {
	return &request{url: url, timeout: timeout}
}

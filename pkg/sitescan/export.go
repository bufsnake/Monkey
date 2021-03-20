package web

var response_header = []string{"Accept-Patch", "Accept-Ranges", "Access-Control-Allow-Origin", "Age", "Allow", "Cache-Control", "Connection", "Content-Disposition", "Content-Encoding", "Content-Language", "Content-Length", "Content-Location", "Content-MD5", "Content-Range", "Content-Security-Policy", "Content-Type", "Date", "ETag", "Expect-CT", "Expires", "Feature-Policy", "Last-Modified", "Link", "Location", "P3P", "Permission-Policy", "Pragma", "Proxy-Authenticate", "Public-Key-Pins", "Referrer-Policy", "Refresh", "Retry-After", "Server", "Set-Cookie", "Status", "Strict-Transport-Security", "Trailer", "Transfer-Encoding", "Upgrade", "Vary", "Via", "WWW-Authenticate", "Warning", "X-Content-Duration", "X-Content-Security-Policy", "X-Content-Type-Options", "X-Frame-Options", "X-Permitted-Cross-Domain-Policies", "X-Powered-By", "X-UA-Compatible", "X-WebKit-CSP", "X-XSS-Protection"}

func NewHttpx(url string, timeout int) *httpx {
	return &httpx{url: url, timeout: timeout}
}

func NewRequest(url string, timeout int) *request {
	return &request{url: url, timeout: timeout}
}

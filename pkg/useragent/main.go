package useragent

import (
	"fmt"
	"math/rand"
)

func RandomUserAgent() string {
	return []func() string{genFirefoxUA, genChromeUA, genOperaUA}[rand.Intn(3)]()
}

var ffVersions = []float32{
	35.0,
	40.0,
	41.0,
	44.0,
	45.0,
	48.0,
	48.0,
	49.0,
	50.0,
	52.0,
	52.0,
	53.0,
	54.0,
	56.0,
	57.0,
	57.0,
	58.0,
	58.0,
	59.0,
	6.0,
	60.0,
	61.0,
	63.0,
}

var chromeVersions = []string{
	"37.0.2062.124",
	"40.0.2214.93",
	"41.0.2228.0",
	"49.0.2623.112",
	"55.0.2883.87",
	"56.0.2924.87",
	"57.0.2987.133",
	"61.0.3163.100",
	"63.0.3239.132",
	"64.0.3282.0",
	"65.0.3325.146",
	"68.0.3440.106",
	"69.0.3497.100",
	"70.0.3538.102",
	"74.0.3729.169",
}

var operaVersions = []string{
	"2.7.62 Version/11.00",
	"2.2.15 Version/10.10",
	"2.9.168 Version/11.50",
	"2.2.15 Version/10.00",
	"2.8.131 Version/11.11",
	"2.5.24 Version/10.54",
}

var osStrings = []string{
	"Macintosh; Intel Mac OS X 10_10",
	"Windows NT 10.0",
	"Windows NT 5.1",
	"Windows NT 6.1; WOW64",
	"Windows NT 6.1; Win64; x64",
	"X11; Linux x86_64",
}

func genFirefoxUA() string {
	version := ffVersions[rand.Intn(len(ffVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s; rv:%.1f) Gecko/20100101 Firefox/%.1f", os, version, version)
}

func genChromeUA() string {
	version := chromeVersions[rand.Intn(len(chromeVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Mozilla/5.0 (%s) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%s Safari/537.36", os, version)
}

func genOperaUA() string {
	version := operaVersions[rand.Intn(len(operaVersions))]
	os := osStrings[rand.Intn(len(osStrings))]
	return fmt.Sprintf("Opera/9.80 (%s; U; en) Presto/%s", os, version)
}

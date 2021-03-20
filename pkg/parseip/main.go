package parseip

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// 支持常见的ip格式
// 192.168.113.159
// 192.168.113.159-254
// 192.168.113.159-192.168.113.254
// 192.168.113.0/24
// 191.168.113.159-192.168.114.254
// 192.167.113.159-192.168.114.254
// 192.168.113.159-192.168.114.254
func ParseIP(ip string) ([]string, error) {
	if strings.Contains(ip, "-") {
		if len(strings.Split(ip, "-")[1]) <= 3 {
			return MultipleIP(ip)
		} else {
			return MultipleIP2(ip)
		}
	} else if strings.Contains(ip, "/") {
		return MultipleIP3(ip)
	} else {
		return SingleIP(ip)
	}
}

// 192.168.113.159
func SingleIP(ip string) ([]string, error) {
	var ret = make([]string, 0)
	for _, val := range strings.Split(ip, ".") {
		ips, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.New(ip + " " + err.Error() + " ip parse error")
		}
		if ips > 255 {
			return nil, errors.New(ip + " ip parse error")
		}
	}
	ret = append(ret, ip)
	return ret, nil
}

// 192.168.113.159-255
func MultipleIP(ips string) ([]string, error) {
	var ret = make([]string, 0)
	host := strings.Split(ips, "-")
	ip := host[0]
	start, err := strconv.Atoi(strings.Split(ip, ".")[3])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	end, err := strconv.Atoi(host[1])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	if start > end {
		return nil, errors.New(ips + " ip parse error")
	}
	if start < 0 {
		start = 0
	}
	if end > 255 {
		end = 255
	}
	for i := start; i <= end; i++ {
		temp := strings.Split(ip, ".")[:3]
		ret = append(ret, temp[0]+"."+temp[1]+"."+temp[2]+"."+strconv.Itoa(i))
	}
	return ret, nil
}

// 192.168.113.159-192.168.113.254
func MultipleIP2(ips string) ([]string, error) {
	var ret = make([]string, 0)
	start := strings.Split(strings.Split(ips, "-")[0], ".")
	end := strings.Split(strings.Split(ips, "-")[1], ".")
	for i := 0; i < 3; i++ {
		if start[i] != end[i] {
			return MultipleIP4(ips)
		}
	}
	temp1, err := strconv.Atoi(start[3])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	temp2, err := strconv.Atoi(end[3])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	if temp1 > temp2 {
		return nil, errors.New(ips + " ip parse error")
	}
	if temp1 < 0 {
		temp1 = 0
	}
	if temp1 >= 255 {
		temp1 = 254
	}
	if temp2 < 0 {
		temp2 = 0
	}
	if temp2 >= 255 {
		temp2 = 254
	}
	for i := temp1; i <= temp2; i++ {
		ret = append(ret, start[0]+"."+start[1]+"."+start[2]+"."+strconv.Itoa(i))
	}
	return ret, nil
}

// 192.168.113.0/24
func MultipleIP3(ips string) ([]string, error) {
	var ret = make([]string, 0)
	host := strings.Split(ips, "/")[0]
	mask, err := strconv.Atoi(strings.Split(ips, "/")[1])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	if len(strings.Split(host, ".")) != 4 {
		return nil, errors.New(ips + " ip parse error")
	}
	a, err := strconv.Atoi(strings.Split(host, ".")[0])
	b, err := strconv.Atoi(strings.Split(host, ".")[1])
	c, err := strconv.Atoi(strings.Split(host, ".")[2])
	d, err := strconv.Atoi(strings.Split(host, ".")[3])
	if err != nil {
		return nil, errors.New(ips + " ip parse error")
	}
	ipbin := fmt.Sprintf("%08s", strconv.FormatInt(int64(a), 2)) +
		fmt.Sprintf("%08s", strconv.FormatInt(int64(b), 2)) +
		fmt.Sprintf("%08s", strconv.FormatInt(int64(c), 2)) +
		fmt.Sprintf("%08s", strconv.FormatInt(int64(d), 2))

	start := ipbin[:mask]
	end := ipbin[:mask]
	for i := 0; i < len(ipbin)-mask; i++ {
		start += "0"
		end += "1"
	}
	start1, err := strconv.ParseUint(start, 2, 32)
	end2, err := strconv.ParseUint(end, 2, 32)
	for i := start1; i <= end2; i++ {
		temp := fmt.Sprintf("%08s", strconv.FormatInt(int64(i), 16))
		e, err := strconv.ParseUint(temp[0:0+2], 16, 8)
		f, err := strconv.ParseUint(temp[2:2+2], 16, 8)
		g, err := strconv.ParseUint(temp[4:4+2], 16, 8)
		h, err := strconv.ParseUint(temp[6:6+2], 16, 8)
		if err != nil {
			return nil, errors.New(ips + " " + err.Error() + " ip parse error")
		}
		ret = append(ret, strconv.Itoa(int(e))+"."+strconv.Itoa(int(f))+"."+strconv.Itoa(int(g))+"."+strconv.Itoa(int(h)))
	}
	return ret, nil
}

// 191.168.113.159-192.168.114.254
// 192.167.113.159-192.168.114.254
// 192.168.113.159-192.168.114.254
func MultipleIP4(ips string) ([]string, error) {
	var ret = make([]string, 0)
	start := strings.Split(strings.Split(ips, "-")[0], ".")
	end := strings.Split(strings.Split(ips, "-")[1], ".")
	var i = 0
	for i = 0; i < 3; i++ {
		if start[i] != end[i] {
			break
		}
	}
	temp1, err := strconv.Atoi(start[i])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error())
	}
	temp2, err := strconv.Atoi(end[i])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error())
	}
	if temp1 > temp2 {
		return nil, errors.New(ips + " parse error")
	}
	a, err := strconv.Atoi(start[0])
	b, err := strconv.Atoi(start[1])
	c, err := strconv.Atoi(start[2])
	d, err := strconv.Atoi(start[3])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	if a > 255 {
		a = 255
	}
	if b > 255 {
		b = 255
	}
	if c > 255 {
		c = 255
	}
	if d > 255 {
		d = 255
	}
	one := fmt.Sprintf("%02s", strconv.FormatInt(int64(a), 16)) +
		fmt.Sprintf("%02s", strconv.FormatInt(int64(b), 16)) +
		fmt.Sprintf("%02s", strconv.FormatInt(int64(c), 16)) +
		fmt.Sprintf("%02s", strconv.FormatInt(int64(d), 16))
	first, err := strconv.ParseUint(one, 16, 32)
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	a, err = strconv.Atoi(end[0])
	b, err = strconv.Atoi(end[1])
	c, err = strconv.Atoi(end[2])
	d, err = strconv.Atoi(end[3])
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	if a > 255 {
		a = 255
	}
	if b > 255 {
		b = 255
	}
	if c > 255 {
		c = 255
	}
	if d > 255 {
		d = 255
	}
	one = fmt.Sprintf("%02s", strconv.FormatInt(int64(a), 16)) +
		fmt.Sprintf("%02s", strconv.FormatInt(int64(b), 16)) +
		fmt.Sprintf("%02s", strconv.FormatInt(int64(c), 16)) +
		fmt.Sprintf("%02s", strconv.FormatInt(int64(d), 16))
	second, err := strconv.ParseUint(one, 16, 32)
	if err != nil {
		return nil, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	for i := first; i <= second; i++ {
		temp := fmt.Sprintf("%08s", strconv.FormatInt(int64(i), 16))
		e, err := strconv.ParseUint(temp[0:0+2], 16, 8)
		f, err := strconv.ParseUint(temp[2:2+2], 16, 8)
		g, err := strconv.ParseUint(temp[4:4+2], 16, 8)
		h, err := strconv.ParseUint(temp[6:6+2], 16, 8)
		if err != nil {
			return nil, errors.New(ips + " " + err.Error() + " ip parse error")
		}
		ret = append(ret, strconv.Itoa(int(e))+"."+strconv.Itoa(int(f))+"."+strconv.Itoa(int(g))+"."+strconv.Itoa(int(h)))
	}
	return ret, nil
}

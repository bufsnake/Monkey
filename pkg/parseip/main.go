package parseip

import (
	"errors"
	"fmt"
	"net"
	"sort"
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
func ParseIP(ip string) (startx uint32, endx uint32, err error) {
	if strings.Contains(ip, "-") {
		if len(strings.Split(ip, "-")[1]) <= 3 {
			return multipleip(ip)
		} else {
			return multipleip2(ip)
		}
	} else if strings.Contains(ip, "/") {
		return multipleip3(ip)
	} else {
		return singleip(ip)
	}
}

// 192.168.113.159
func singleip(ip string) (startx uint32, endx uint32, err error) {
	for _, val := range strings.Split(ip, ".") {
		ips, err := strconv.Atoi(val)
		if err != nil {
			return 0, 0, errors.New(ip + " " + err.Error() + " ip parse error")
		}
		if ips > 255 {
			return 0, 0, errors.New(ip + " ip parse error")
		}
	}
	return ip2UInt32(ip), ip2UInt32(ip), nil
}

// 192.168.113.159-255
func multipleip(ips string) (startx uint32, endx uint32, err error) {
	host := strings.Split(ips, "-")
	ip := host[0]
	start, err := strconv.Atoi(strings.Split(ip, ".")[3])
	if err != nil {
		return 0, 0, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	end, err := strconv.Atoi(host[1])
	if err != nil {
		return 0, 0, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	if start > end {
		return 0, 0, errors.New(ips + " ip parse error")
	}
	if start < 0 {
		start = 0
	}
	if end > 255 {
		end = 255
	}
	temp := strings.Split(ip, ".")
	return ip2UInt32(temp[0] + "." + temp[1] + "." + temp[2] + "." + strconv.Itoa(start)), ip2UInt32(temp[0] + "." + temp[1] + "." + temp[2] + "." + strconv.Itoa(end)), nil
}

// 192.168.113.159-192.168.113.254
func multipleip2(ips string) (startx uint32, endx uint32, err error) {
	start := ip2UInt32(strings.Split(ips, "-")[0])
	end := ip2UInt32(strings.Split(ips, "-")[1])
	if start > end {
		return 0, 0, errors.New(ips + " error")
	}
	return start, end, nil
}

// 192.168.113.0/24
func multipleip3(ips string) (startx uint32, endx uint32, err error) {
	host := strings.Split(ips, "/")[0]
	mask, err := strconv.Atoi(strings.Split(ips, "/")[1])
	if err != nil {
		return 0, 0, errors.New(ips + " " + err.Error() + " ip parse error")
	}
	if len(strings.Split(host, ".")) != 4 {
		return 0, 0, errors.New(ips + " ip parse error")
	}
	a, err := strconv.Atoi(strings.Split(host, ".")[0])
	b, err := strconv.Atoi(strings.Split(host, ".")[1])
	c, err := strconv.Atoi(strings.Split(host, ".")[2])
	d, err := strconv.Atoi(strings.Split(host, ".")[3])
	if err != nil {
		return 0, 0, errors.New(ips + " ip parse error")
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
	if err != nil {
		return 0, 0, errors.New(ips + " ip parse error: " + err.Error())
	}
	end2, err := strconv.ParseUint(end, 2, 32)
	if err != nil {
		return 0, 0, errors.New(ips + " ip parse error: " + err.Error())
	}
	return uint32(start1), uint32(end2), nil
}

func ip2UInt32(ipnr string) uint32 {
	bits := strings.Split(ipnr, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum uint32
	sum += uint32(b0) << 24
	sum += uint32(b1) << 16
	sum += uint32(b2) << 8
	sum += uint32(b3)
	return sum
}

func UInt32ToIP(intIP uint32) string {
	var bytes [4]byte
	bytes[0] = byte(intIP & 0xFF)
	bytes[1] = byte((intIP >> 8) & 0xFF)
	bytes[2] = byte((intIP >> 16) & 0xFF)
	bytes[3] = byte((intIP >> 24) & 0xFF)

	return net.IPv4(bytes[3], bytes[2], bytes[1], bytes[0]).String()
}

// 去重
// 本身为啥就为啥
// 本身为反向则清空该点
// 本身为center则设为目前点
// 第一个一定是right
// 判断第二个 一直找，知道找到下一个right的前一个left
func DeDuplication(allip [][2]uint32) [][2]uint32 {
	var results [][2]uint32
	linear := make(map[uint32]string) // left center right
	for i := 0; i < len(allip); i++ {
		if allip[i][0] == allip[i][1] {
			if _, b := linear[allip[i][0]]; !b {
				linear[allip[i][0]] = "center"
			}
			continue
		}
		if _, b := linear[allip[i][0]]; !b {
			linear[allip[i][0]] = "right"
		} else {
			if linear[allip[i][0]] == "left" {
				delete(linear, allip[i][0])
			} else if linear[allip[i][0]] == "center" {
				linear[allip[i][0]] = "right"
			}
		}
		if _, b := linear[allip[i][1]]; !b {
			linear[allip[i][1]] = "left"
		} else {
			if linear[allip[i][1]] == "right" {
				delete(linear, allip[i][0])
			} else if linear[allip[i][1]] == "center" {
				linear[allip[i][0]] = "left"
			}
		}
	}
	key := make([]int, 0)
	for xxx, _ := range linear {
		key = append(key, int(xxx))
	}
	sort.Ints(key)
	for i := 0; i < len(key); i++ {
		var start, end uint32
		start = uint32(key[i])
		for j := i + 1; j < len(key); j++ {
			if (linear[uint32(key[j])] == "left" && j+1 == len(linear)) || (linear[uint32(key[j])] == "left" && linear[uint32(key[j+1])] == "right") {
				end = uint32(key[j])
				i = j
				break
			}
		}
		results = append(results, [2]uint32{start, end})
	}
	return results
}

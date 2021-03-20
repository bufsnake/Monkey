package log

import (
	"fmt"
	. "github.com/logrusorgru/aurora"
	"math"
	"strings"
	"sync"
	"time"
)

var totalcountbak float64    // 总数
var totalcount float64       // 总数
var increase float64         // 递增
var schedule_lock sync.Mutex // 进度锁
var ipalive map[string]int   // IP存活数量
var portnum int              // 端口数量
var errnum int               // 失败数量
var timestart time.Time      // 开始时间

func Bar() {
	for {
		percentage := math.Trunc(((increase/totalcount)*100)*1e2) * 1e-2 // 更新进度条
		fmt.Printf("\r%s %.0f  %s %.2f%%  %s %d  %s %d  %s %d  %s %.2fs               ", BrightWhite("AllIP:").String(), totalcountbak, BrightWhite("Percentage:").String(), percentage, BrightWhite("IP:").String(), len(ipalive), BrightWhite("Port:").String(), portnum, BrightWhite("Err:").String(), errnum, BrightWhite("Time:").String(), time.Now().Sub(timestart).Seconds())
		time.Sleep(1 * time.Second)
	}
}

// 首先根据IP数量设置totalcount
func SetTotalCount(ipcount int) {
	totalcount = float64(ipcount)
	totalcountbak = float64(ipcount)
	timestart = time.Now()
	ipalive = make(map[string]int)
	portnum = 0
	errnum = 0
}

func UpdateTotalCount(portcount int) {
	schedule_lock.Lock()
	defer schedule_lock.Unlock()
	totalcount += float64(portcount - 1)
	percentage := math.Trunc(((increase/totalcount)*100)*1e2) * 1e-2 // 更新进度条
	fmt.Printf("\r%s %.0f  %s %.2f%%  %s %d  %s %d  %s %d  %s %.2fs               ", BrightWhite("AllIP:").String(), totalcountbak, BrightWhite("Percentage:").String(), percentage, BrightWhite("IP:").String(), len(ipalive), BrightWhite("Port:").String(), portnum, BrightWhite("Err:").String(), errnum, BrightWhite("Time:").String(), time.Now().Sub(timestart).Seconds())
}

// Terminal == IP PORT PROTOCOL VERSION
// CSV      == IP PORT PROTOCOL SCREEN SERVICEFP
// 进度: Percentage IP存活: IPAliveNum 端口数量: PortNum 当前耗时: xxxx
func Println(a ...interface{}) {
	schedule_lock.Lock()
	defer schedule_lock.Unlock()
	increase++
	terminal := []interface{}{}
	for i := 0; i < len(a); i++ {
		if i == 4 {
			break
		}
		if i == 0 {
			terminal = append(terminal, "\r"+color(i, a[i]))
		} else {
			terminal = append(terminal, color(i, a[i]))
		}
	}
	if len(a) > 2 {
		ipalive[color(0, a[0])] = 1
		portnum++
		if len(terminal) == 4 {
			if len(color(3, a[3])) < 30+9 {
				terminal = append(terminal, strings.Repeat(" ", 30))
			}

		}
	} else {
		errnum++
		terminal = append(terminal, strings.Repeat(" ", 30))
	}
	fmt.Println(terminal...)
	if totalcount == 0 {
		return
	}
	percentage := math.Trunc(((increase/totalcount)*100)*1e2) * 1e-2 // 更新进度条
	fmt.Printf("\r%s %.0f  %s %.2f%%  %s %d  %s %d  %s %d  %s %.2fs               ", BrightWhite("AllIP:").String(), totalcountbak, BrightWhite("Percentage:").String(), percentage, BrightWhite("IP:").String(), len(ipalive), BrightWhite("Port:").String(), portnum, BrightWhite("Err:").String(), errnum, BrightWhite("Time:").String(), time.Now().Sub(timestart).Seconds())
}

func color(i int, data interface{}) string {
	switch i {
	case 0: // IP
		return fmt.Sprintf("%-25s", BrightGreen(data).String())
	case 1: // PORT
		return fmt.Sprintf("%-14s", BrightCyan(data).String())
	case 2: // PROTOCOL
		return fmt.Sprintf("%-34s", BrightMagenta(data).String())
	case 3: // VERSION
		return fmt.Sprintf("%s", BrightYellow(data).String())
	default:
		return BrightWhite(data).String()
	}
}

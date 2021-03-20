package main

import (
	"flag"
	"fmt"
	"github.com/bufsnake/Monkey/config"
	"github.com/bufsnake/Monkey/internal/core"
	"github.com/bufsnake/Monkey/internal/other"
	"github.com/bufsnake/Monkey/pkg/log"
	"github.com/bufsnake/Monkey/pkg/parseip"
	"github.com/logrusorgru/aurora"
	"io/ioutil"
	"strings"
)

func main() {
	conf := config.Config{}
	flag.StringVar(&conf.IP, "t", "", "specify target ip")
	flag.StringVar(&conf.File, "f", "", "specify target ip from file")
	flag.IntVar(&conf.Thread, "r", 6, "specify scan rate")
	flag.IntVar(&conf.Timeout, "w", 10, "web request timeout")
	flag.StringVar(&conf.NmapsV, "n", "2", "nmap version intensity,optional 0~9")
	flag.StringVar(&conf.Port, "p", "0-65535", "specify scan ports")
	flag.StringVar(&conf.MasscanRate, "m", "1000", "masscan rate")
	flag.Parse()
	if conf.IP == "" && conf.File == "" {
		flag.Usage()
		return
	}
	allip := []string{}
	if conf.IP != "" {
		ip, err := parseip.ParseIP(strings.Trim(conf.IP, " "))
		if err != nil {
			log.Println(err)
			return
		}
		for i := 0; i < len(ip); i++ {
			if !other.Exist(allip, ip[i]) {
				allip = append(allip, ip[i])
			}
		}
	} else if conf.File != "" {
		file, err := ioutil.ReadFile(conf.File)
		if err != nil {
			log.Println(err)
			return
		}
		split := strings.Split(string(file), "\n")
		for i := 0; i < len(split); i++ {
			trim := strings.Trim(split[i], "\r")
			if len(trim) == 0 {
				continue
			}
			if len(trim) >= 7 {
				ip, err := parseip.ParseIP(trim)
				if err != nil {
					log.Println(err)
					return
				}
				for i := 0; i < len(ip); i++ {
					if !other.Exist(allip, ip[i]) {
						allip = append(allip, ip[i])
					}
				}
			} else {
				log.Println("error data", trim)
				return
			}
		}
	}
	log.SetTotalCount(len(allip))
	xcore := core.NewCore(conf, allip)
	xcore.Start()
	fmt.Println()
	for i := 0; i < len(config.Urldatas); i++ {
		if i == 0 {
			fmt.Println("Urls Results:")
		}
		fmt.Println("["+aurora.BrightGreen(config.Urldatas[i][3]).String()+"]", "["+aurora.BrightMagenta(config.Urldatas[i][0]).String()+"]", "["+aurora.BrightWhite(config.Urldatas[i][1]).String()+"]", "["+aurora.BrightCyan(config.Urldatas[i][2]).String()+"]")

	}
}

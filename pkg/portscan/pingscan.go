package findport

import (
	"errors"
	"github.com/go-ping/ping"
	"time"
)

type message struct {
	success bool
	err     error
}

func pingscan(ip string) (bool, error) {
	messages := make(chan message)
	for i := 0; i < 3; i++ {
		go func() {
			pinger, err := ping.NewPinger(ip)
			if err != nil {
				messages <- message{success: false, err: err}
				return
			}
			pinger.Timeout = 1 * time.Second
			pinger.Count = 1
			err = pinger.Run()
			if err != nil {
				messages <- message{success: false, err: err}
				return
			}
			stats := pinger.Statistics()
			if stats.PacketsRecv > 0 {
				messages <- message{success: true, err: err}
				return
			}
			messages <- message{success: false, err: errors.New("Not Alive")}
		}()
	}
	for i := 0; i < 3; i++ {
		temp := <-messages
		if temp.success {
			return temp.success, temp.err
		}
	}
	return false, errors.New("Not Alive")
}

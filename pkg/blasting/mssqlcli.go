package blasting

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type mssql struct {
	ip   string
	port string
}

func (m *mssql) Info() string {
	return "weak"
}

func (m *mssql) Connect() (string, error) {
	usernames := []string{"sa"}
	passwords := []string{"admin", "%user%", "%user%123", "%user%1234", "%user%123456", "%user%12345", "%user%@123", "%user%@123456", "%user%@12345", "%user%#123", "%user%#123456", "%user%#12345", "%user%_123", "%user%_123456", "%user%_12345", "%user%123!@#", "%user%!@#$", "%user%!@#", "%user%~!@", "%user%!@#123", "qweasdzxc", "%user%2017", "%user%2016", "%user%2015", "%user%@2017", "%user%@2016", "%user%@2015", "Passw0rd", "qweasdzxc", "admin123", "admin888", "administrator", "administrator123", "sa123", "ftp", "ftppass", "123456", "password", "12345", "1234", "sa", "123", "qwerty", "test", "1q2w3e4r", "1qaz2wsx", "qazwsx", "123qwe", "123qaz", "0000", "oracle", "1234567", "123456qwerty", "password123", "12345678", "1q2w3e", "abc123", "okmnji", "test123", "123456789", "q1w2e3r4", "sqlpass", "sql123", "sqlserver", "web", ""}
	wait := sync.WaitGroup{}
	messages := make(chan message, 600)
	fin := make(chan message)
	for i := 0; i < 5; i++ {
		wait.Add(1)
		go m.check(&wait, messages, fin)
	}
	for _, u := range usernames {
		for _, p := range passwords {
			messages <- message{
				user: u,
				pass: p,
			}
		}
	}
	go func() {
		close(messages)
		wait.Wait()
		fin <- message{user: "error"}
	}()
	select {
	case <-time.After(5 * time.Minute):
		return "", errors.New("mssql weak password test timeout")
	case mess := <-fin:
		if mess.user == "error" {
			return "", errors.New("mssql weak password test finish,but no password found")
		}
		return "mssql://" + mess.user + ":" + mess.pass + "@" + m.ip + ":" + m.port, nil
	}
}

func (m *mssql) check(wg *sync.WaitGroup, messages, fin chan message) {
	defer wg.Done()
	for message_ := range messages {
		message_.pass = strings.ReplaceAll(message_.pass, "%user%", message_.user)
		db, err := sql.Open("mssql", fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=%s", m.ip, m.port, message_.user, message_.pass, "master"))
		if err != nil {
			fmt.Println(err)
			continue
		}
		err = db.Ping()
		if err != nil {
			fmt.Println(err)
			db.Close()
			continue
		}
		db.Close()
		fin <- message_
	}
}

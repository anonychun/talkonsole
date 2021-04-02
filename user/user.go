package user

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/anonychun/talkonsole/constant"
)

type User struct {
	Conn net.Conn

	Name string
	Room string
}

func (u *User) Scan(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	scanner.Scan()
	return scanner.Text()
}

func (u *User) SendMessage(message string) error {
	_, err := u.Conn.Write([]byte(message))
	return err
}

func (u *User) SendMessageToRoom(message string) error {
	return u.SendMessage(fmt.Sprintf("%s >> %s", u.Name, message))
}

func (u *User) AcceptMessage() (<-chan string, <-chan error) {
	message := make(chan string)
	errchan := make(chan error)

	go func() {
		defer close(message)
		defer close(errchan)

		for {
			data := make([]byte, constant.BUFFER_SIZE)
			_, err := u.Conn.Read(data)
			if err != nil {
				errchan <- err
				return
			}
			message <- string(data)
		}
	}()

	return message, errchan
}

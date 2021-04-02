package client

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"

	"github.com/anonychun/talkonsole/user"
)

type client struct {
	conn net.Conn
}

func Join(host string, port int, name, room string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	c := &client{conn}

	u := &user.User{
		Name: name,
		Room: room,
	}

	return c.handle(u)
}

func (c *client) handle(u *user.User) error {
	defer c.conn.Close()

	err := c.register(u)
	if err != nil {
		return err
	}
	u.Conn = c.conn

	go func() {
		for {
			message := u.Scan(os.Stdin)
			err = u.SendMessageToRoom(message)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}
	}()

	messages, errchan := u.AcceptMessage()
	for {
		select {
		case message := <-messages:
			fmt.Println(message)
		case <-errchan:
			return nil
		}
	}
}

func (c *client) register(u *user.User) error {
	return gob.NewEncoder(c.conn).Encode(u)
}

package server

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/anonychun/talkonsole/user"
)

type server struct {
	rooms map[string][]*user.User
}

func Start(port int) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	defer listener.Close()

	s := &server{
		rooms: make(map[string][]*user.User),
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		u := &user.User{}
		gob.NewDecoder(conn).Decode(u)
		u.Conn = conn

		go s.handle(u)
	}
}

func (s *server) handle(u *user.User) {
	defer u.Conn.Close()

	if s.checkUserExist(u) {
		u.Conn.Write([]byte("user already exist"))
		return
	}

	err := s.addUserToRoom(u)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	messages, errchan := u.AcceptMessage()
	for {
		select {
		case message := <-messages:
			s.broadcastToRoom(u.Room, message)
		case <-errchan:
			s.removeUserFromRoom(u)
			return
		}
	}
}

func (s *server) checkUserExist(u *user.User) bool {
	for _, v := range s.rooms[u.Room] {
		if v.Name == u.Name {
			return true
		}
	}

	return false
}

func (s *server) broadcastToRoom(room, message string) error {
	for _, u := range s.rooms[room] {
		go u.SendMessage(message)
	}

	return nil
}

func (s *server) addUserToRoom(u *user.User) error {
	err := s.notifyClientJoin(u)
	if err != nil {
		return err
	}

	s.rooms[u.Room] = append(s.rooms[u.Room], u)
	return nil
}

func (s *server) removeUserFromRoom(u *user.User) error {
	index := 0
	for i, v := range s.rooms[u.Room] {
		if v.Name == u.Name {
			index = i
			break
		}
	}
	s.rooms[u.Room] = append(s.rooms[u.Room][:index], s.rooms[u.Room][index+1:]...)

	return s.notifyClientLeft(u)
}

func (s *server) notifyClientJoin(u *user.User) error {
	message := fmt.Sprintf("[+] %s join %s", u.Name, u.Room)
	fmt.Println(message)
	return s.broadcastToRoom(u.Room, message)
}

func (s *server) notifyClientLeft(u *user.User) error {
	message := fmt.Sprintf("[-] %s left %s", u.Name, u.Room)
	fmt.Println(message)
	return s.broadcastToRoom(u.Room, message)
}

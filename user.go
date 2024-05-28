// user define
package main

import (
	"fmt"
	"net"
)

type User struct {
	Name   string
	Ip     string
	Port   int
	Ch     chan string
	Conne  net.Conn
	Server *MudServer
}

func NewUser(name string, conne net.Conn, server *MudServer) *User {
	fmt.Println("usr creating...")
	var user = User{Name: name}
	user.Ch = make(chan string)
	user.Conne = conne
	user.Server = server
	fmt.Println("usr watching channel...")
	//很重要！必须有不同的goruing 来分别对channel进行读写
	go user.Watch()
	server.BroadCast("system:" + user.Name + "is online...")
	return &user
}

func (me *User) Send(msg string) {
	fmt.Println("user send msg", msg)
	me.Conne.Write([]byte(msg + "\n"))

}

func (me *User) Watch() {
	for {
		msg := <-me.Ch
		if msg == "exit" {
			break
		} else {
			me.Send(msg)
		}

	}
}

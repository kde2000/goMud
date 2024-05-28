// user define
package main

import (
	"fmt"
	"net"
)

type User struct {
	Name  string
	Ip    string
	Port  int
	Ch    chan string
	Conne net.Conn
}

func NewUser(name string, conne net.Conn) *User {
	fmt.Println("usr creating...")
	var user = User{Name: name}
	user.Ch = make(chan string)
	user.Conne = conne
	return &user
}

func (me *User) Send(msg string) {
	fmt.Println("user send msg", msg)

}

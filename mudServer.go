package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

//场景，每个场景有个发言板

type Scean struct {
	Location  string
	Billboard chan string
}

type UserScean struct {
	User  User
	Scean Scean
}

var GameWord = Scean{Location: "lobby", Billboard: make(chan string)}
var lock sync.RWMutex

// mud服务器，带玩家和场景
type MudServer struct {
	Ip     string
	Port   int
	Family map[string]UserScean
}

// 初始化服务器
func NewMudServer(ip string, port int) *MudServer {
	fmt.Println("create new mudserver")
	ns := new(MudServer)
	ns.Ip = ip
	ns.Port = port
	ns.Family = make(map[string]UserScean)
	return ns

}

// 服务器启动函数
func (me *MudServer) Start() {

	fmt.Println("Mud server init...")
	//listen 端口
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", me.Ip, me.Port))
	if err != nil {
		fmt.Println("listen problem...")
		return
	}
	defer listener.Close()
	for {
		//accept 连接
		fmt.Println("server listening...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept problem...")
			break
		}
		//登记用户
		fmt.Println("server accepting...", conn.RemoteAddr().String())
		newbie := *NewUser(conn.RemoteAddr().String(), conn, me)
		lock.Lock()
		me.Family[conn.RemoteAddr().String()] = UserScean{User: newbie, Scean: GameWord}
		lock.Unlock()
		//启动新的协程处理server
		go me.Handle(conn)
	}
}

func (me *MudServer) Handle(conn net.Conn) {
	fmt.Println("User in:", me.Family)
	defer conn.Close()
	//很重要。匿名函数可以使用主函数的变量
	//一个单独的协程来处理读入数据
	go func() {
		var buff = make([]byte, 4096)
		for {
			n, err := conn.Read(buff)
			if err != nil && err != io.EOF {
				fmt.Println("read error")
				return
			}
			if n > 0 {
				me.BroadCast(string(buff[:n-1]))
			}
		}

	}()
	//很重要。机制不清楚。
	select {}

}
func (me *MudServer) BroadCast(msg string) {
	lock.Lock()
	for _, v := range me.Family {
		v.User.Ch <- msg
	}
	lock.Unlock()
}

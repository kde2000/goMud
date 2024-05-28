package main

import (
	"fmt"
	"net"
)

//场景，每个场景有个发言板

type Scean struct {
	Location  string
	Billboard chan string
}

var GameWord = Scean{Location: "lobby", Billboard: make(chan string)}

// mud服务器，带玩家和场景
type MudServer struct {
	Ip     string
	Port   int
	Family map[string]Scean
}

// 初始化服务器
func NewMudServer(ip string, port int) *MudServer {
	fmt.Println("create new mudserver")
	ns := new(MudServer)
	ns.Ip = ip
	ns.Port = port
	ns.Family = make(map[string]Scean)
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
	for {
		//accept 连接
		fmt.Println("server listening...")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept problem...")
			return
		}
		//登记用户
		fmt.Println("server accepting...", conn.RemoteAddr().String())
		me.Family[conn.RemoteAddr().String()] = GameWord
		//启动新的协程
		go me.Handle(conn)
	}
}

func (me *MudServer) Handle(conn net.Conn) {
	fmt.Println("User in:", me.Family)

}

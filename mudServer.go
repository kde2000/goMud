package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

//场景，每个场景有个发言板

type Scean struct {
	Location  string
	Billboard chan string
}

type UserScean struct {
	User  *User
	Scean *Scean
}

var GameWord = &Scean{Location: "lobby", Billboard: make(chan string)}
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
		newbie := NewUser(conn.RemoteAddr().String(), conn, me)
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
				inString := strings.Split(string(buff[:n-1]), "|")
				comd := strings.Trim(inString[0], " ")
				fmt.Println("comd in:", comd)
				switch comd {
				case "all":
					me.BroadCast(string(buff[:n-1]))
				case "who", "WHO":
					me.Who()
				case "exit", "EXIT":
					me.BroadCast(conn.RemoteAddr().String() + " offline...")
					return
				default:
					if len(inString) > 1 {
						msgtxt := inString[1]
						me.SendUserMessage(comd, msgtxt)
					} else {
						me.SendUserMessage(conn.RemoteAddr().String(), "格式问题：用户名|信息。命令：all，who，exit")
					}

				}

			}
		}

	}()
	select {}

}
func (me *MudServer) BroadCast(msg string) {
	lock.Lock()
	for _, v := range me.Family {
		v.User.Ch <- msg
	}
	lock.Unlock()
}
func (me *MudServer) Who() {
	lock.Lock()
	for _, v := range me.Family {
		me.BroadCast(v.User.Name)
	}
	lock.Unlock()
}
func (me *MudServer) SendUserMessage(username, msg string) {
	lock.Lock()
	for _, v := range me.Family {
		if v.User.Name == username {
			v.User.Ch <- msg
		}

	}
	lock.Unlock()
}

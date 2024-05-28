package main

import "fmt"

func main() {
	fmt.Println("server created...")
	var gameServer = NewMudServer("192.168.2.212", 9999)
	gameServer.Start()
}

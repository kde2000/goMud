package main

import "fmt"

func main() {
	fmt.Println("server created...")
	var gameServer = NewMudServer("127.0.0.1", 9999)
	gameServer.Start()
}

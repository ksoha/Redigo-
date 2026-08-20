package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	//Open a TCP port and start listening for incoming connections

	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Println("Redigo listening on port 6380")

	//loop forever waiting for incoming clients to connect
	//allow() itself blocks so the loop doesnt consume one million times
	//each time someone connects, add them to handleconnection in a new goroutine

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close() //close the co
	fmt.Println("New client connected from", conn.RemoteAddr())

	//just prove connection works
	conn.Write([]byte("Welcome to Redigo!\n"))
}

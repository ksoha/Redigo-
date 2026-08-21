package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/ksoha/redigo/internal/resp"
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

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	//keep reading commands froms this client until thy disconnect
	//or an error occurs
	for {
		//read a commanf from the client
		args, err := resp.ReadCommand(reader)
		if err != nil {
			//client disconnect or send malformed info
			//stop handling this connection
			fmt.Println("Error reading command:", conn.RemoteAddr(), "reason:", err)
			return
		}

		fmt.Println("recieved command: %v\n", args)
		if err := resp.WriteSimpleString(writer, "OK"); err != nil {
			fmt.Println("write error:", err)
			return
		}
	}
}

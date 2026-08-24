package main

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/ksoha/redigo/internal/command"
	"github.com/ksoha/redigo/internal/resp"
	"github.com/ksoha/redigo/internal/store"
)

func main() {
	//Open a TCP port and start listening for incoming connections

	listener, err := net.Listen("tcp", ":6380")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	fmt.Println("Redigo listening on port 6380")

	//create one shared store for all clients to use
	//every client will use the same store
	s := store.New()

	//loop forever waiting for incoming clients to connect
	//allow() itself blocks so the loop doesnt consume one million times
	//each time someone connects, add them to handleconnection in a new goroutine
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)

			continue
		}

		go handleConnection(conn, s)
	}
}

func handleConnection(conn net.Conn, s *store.Store) {
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

		fmt.Printf("recieved command: %v\n", args)

		//dispatch the command to the store and write the response back to the client \
		if err := command.Dispatch(args, s, writer); err != nil {
			fmt.Println("dispatch error", err)
			return
		}
	}
}

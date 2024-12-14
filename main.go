package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	s "tcp/pkg/server"
	"github.com/pebbe/zmq4"
)

func main() {
	fmt.Println("Starting the server...")

	host := "0.0.0.0"
	port := "6969"

	server := *s.NewServer(host, port)
	fmt.Println("Server created")

	err := server.OpenSocket()
	if err != nil {
		fmt.Printf("Failed to open socket: %v\n", err)
		os.Exit(1)
	}
	defer server.Socket.Close()

	fmt.Printf("Server is listening on %s:%s\n", host, port)

	// zmq stuff

	publisher, err_create_sock := zmq4.NewSocket(zmq4.PUB)
	if err_create_sock != nil {
		fmt.Println("Error while creating ZMQ socket", err)
		return
	}
	defer publisher.Close()

	err_create_pub := publisher.Bind("tcp://localhost:6970")
	if err_create_pub != nil {
		fmt.Println("Error while creating publisher", err)
	}

	for {
		err := server.OpenConnection()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		fmt.Println("New client connected")
		go handleClient(server.Connection, publisher)
	}
}

func handleClient(conn net.Conn, publisher *zmq4.Socket) {
	defer conn.Close()
	fmt.Printf("Handling connection from %s\n", conn.RemoteAddr())

	buffer := make([]byte, 1024)
	for {

		length, err := conn.Read(buffer)

		if err != nil {
			fmt.Printf("Client %s disconnected or error reading: %v\n", conn.RemoteAddr(), err)
			break
		}

		msg := string(buffer[:length])

		var data interface{}
		json.Unmarshal([]byte(msg), &data)

		prettyJSON, err := json.MarshalIndent(data, "", "  ")

		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		data, err_send := publisher.Send(msg, zmq4.DONTWAIT)
		if err_send != nil {
			fmt.Println("Error when sending", err_send)
		}

		fmt.Printf("RECV from %s:\n%s\n", conn.RemoteAddr(), prettyJSON)

	}
}

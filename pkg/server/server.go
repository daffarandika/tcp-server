package server

import (
	"fmt"
	"net"
)

type Server struct {
	Address     string
	Connection  net.Conn
	Socket    net.Listener
}

func NewServer(host, port string) *Server {
	return &Server{
		Address:     host + ":" + port,
	}
}

func (server *Server) OpenSocket() error {
	listener, err := net.Listen("tcp", server.Address)

	if err != nil {
		return err
	}

	server.Socket = listener

	return nil
}

func (server *Server) OpenConnection() error {
	conn, err := server.Socket.Accept()
	if err != nil {
		fmt.Println("error is here open")
		return err
	}
	server.Connection = conn
	return nil
}

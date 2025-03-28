package server

import (
	"bufio"
	"log"
	"net"

	reply "github.com/Anurag-Raut/smtp/dto"
)

type Server struct {
	port     string
	adddress string
}

type Config struct {
	port    *string
	address *string
}

func NewConfig() Config {
	return Config{}
}
func (c *Config) SetPort(port string) {
	c.port = &port
}
func (c *Config) SetAddr(address string) {
	c.address = &address
}

func NewServer(c Config) *Server {
	addr := "127.0.0.1"
	port := "587"
	if c.address != nil {
		addr = *c.address
	}
	if c.port != nil {
		port = *c.port
	}
	server := Server{
		port:     port,
		adddress: addr,
	}
	return &server
}

func (s *Server) Listen() {
	listner, err := net.Listen("tcp", s.adddress+":"+s.port)
	if err != nil {
		log.Println("Error: ", err.Error())
	}

	for {
		c, err := listner.Accept()
		if err != nil {
			log.Print("Error", err.Error())
		}
		go handleConn(c)
	}
}

func handleConn(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	reply.Greet(writer)

	handleCommands(reader, writer)

}

func greet(w *bufio.Writer) {
	newGreeting := reply.NewGreeting()

	w.Write(newGreeting)
}

func handleCommands(r *bufio.Reader, w *bufio.Writer) {

}

package server

import (
	"bufio"
	"log"
	"net"

	"github.com/Anurag-Raut/smtp/server/io/reader"
	"github.com/Anurag-Raut/smtp/server/session"
)

type Server struct {
	port     string
	adddress string
	listner  net.Listener
	done     chan struct{}
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
		done:     make(chan struct{}),
	}
	return &server
}

func (s *Server) Listen() {
	log.Println("Listening on port", s.port)
	newListner, err := net.Listen("tcp", s.adddress+":"+s.port)
	s.listner = newListner
	if err != nil {
		log.Println("Error: ", err.Error())
	}

	for {
		c, err := s.listner.Accept()
		if err != nil {
			select {
			case <-s.done:
				log.Println("server is shutting down")
				return
			default:
				log.Print("Error", err.Error())
			}
		}
		go handleConn(c)
	}
}

func (s *Server) Close() {
	close(s.done)
	s.listner.Close()
}

func handleConn(conn net.Conn) {
	reader := reader.NewReader(conn)
	writer := bufio.NewWriter(conn)
	session := session.NewSession()
	session.Begin(reader, writer)

}

package session

import (
	"bufio"
	"net"

	"github.com/Anurag-Raut/smtp/client/parser"
)

type Session struct {
	reader *bufio.Reader
	writer *bufio.Writer
	parser
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (s *Session) SendEmail(from string, to []string, body *string) {

}

func (s *Session) Begin() {
	parser := parser.
		reply.getReply()
}

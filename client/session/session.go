package session

import (
	"bufio"
	"net"

	"github.com/Anurag-Raut/smtp/client/dto/command"
	"github.com/Anurag-Raut/smtp/client/dto/reply"
	"github.com/Anurag-Raut/smtp/client/parser"
)

type Session struct {
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewSession(conn net.Conn) *Session {
	return &Session{
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (s *Session) SendEmail(from string, to []string, body *string) {

}

func (s *Session) Begin() error {
	p := parser.NewReplyParser(s.reader)
	greeting, err := reply.GetReply(parser.Greeting, p)
	if err != nil {
		return err
	}
	command.SendEHLO(s.writer)

	return nil
}

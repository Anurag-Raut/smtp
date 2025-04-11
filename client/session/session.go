package session

import (
	"bufio"
	"log"
	"net"
	"net/http"

	"github.com/Anurag-Raut/smtp/client/dto/command"
	"github.com/Anurag-Raut/smtp/client/dto/reply"
	"github.com/Anurag-Raut/smtp/client/parser"
)

type Session struct {
	reader     *bufio.Reader
	writer     *bufio.Writer
	httpWriter http.ResponseWriter
}

func NewSession(conn net.Conn, w http.ResponseWriter) *Session {
	return &Session{
		reader:     bufio.NewReader(conn),
		writer:     bufio.NewWriter(conn),
		httpWriter: w,
	}
}

func (s *Session) SendEmail(from string, to []string, body *string) {

	command.SendMail(s.writer, from)
	command.SendRcpt(s.writer, to[0])
	if body != nil {
		command.SendBody(s.writer, *body)
	}

}

func (s *Session) Begin() error {
	log.Println("wa")
	p := parser.NewReplyParser(s.reader)
	_, err := reply.GetReply(parser.Greeting, p)
	if err != nil {
		return err
	}
	command.SendEHLO(s.writer)
	reply.GetReply(parser.ReplyLine, p)
	return nil
}

package session

import (
	"bufio"
	"net"
	"net/http"

	"github.com/Anurag-Raut/smtp/client/dto/command"
	"github.com/Anurag-Raut/smtp/client/dto/reply"
	"github.com/Anurag-Raut/smtp/client/io/reader"
	"github.com/Anurag-Raut/smtp/client/parser"
	"github.com/Anurag-Raut/smtp/logger"
)

type Session struct {
	smtpConn   net.Conn
	reader     *reader.Reader
	writer     *bufio.Writer
	httpWriter http.ResponseWriter
}

func NewSession(conn net.Conn, w http.ResponseWriter) *Session {
	return &Session{
		reader:     reader.NewReader(conn),
		writer:     bufio.NewWriter(conn),
		httpWriter: w,
		smtpConn:   conn,
	}
}

func (s *Session) SendEmail(from string, to []string, body *string) {

	command.SendMail(s.writer, from)
	command.SendRcpt(s.writer, to[0])
	if body != nil {
		command.SendBody(s.writer, *body)
	}
	command.SendQuit(s.writer)
	s.smtpConn.Close()
}

func (s *Session) Begin() error {
	p := parser.NewReplyParser(s.reader)

	_, err := reply.GetReply(parser.Greeting, p)
	if err != nil {
		return err
	}
	logger.ClientLogger.Println("Parsed Greeting")
	command.SendEHLO(s.writer)
	reply.GetReply(parser.ReplyLine, p)
	return nil
}

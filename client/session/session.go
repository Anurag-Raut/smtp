package session

import (
	"net"
	"net/http"

	"github.com/Anurag-Raut/smtp/client/dto/command"
	"github.com/Anurag-Raut/smtp/client/dto/reply"
	"github.com/Anurag-Raut/smtp/client/io/reader"
	"github.com/Anurag-Raut/smtp/client/io/writer"
	"github.com/Anurag-Raut/smtp/client/parser"
)

type Session struct {
	smtpConn   net.Conn
	reader     *reader.Reader
	writer     *writer.Writer
	httpWriter http.ResponseWriter
}

func NewSession(conn net.Conn, w http.ResponseWriter) *Session {
	return &Session{
		reader:     reader.NewReader(conn),
		writer:     writer.NewWriter(conn),
		httpWriter: w,
		smtpConn:   conn,
	}
}

func (s *Session) SendEmail(from string, to []string, body *string) {

	p := parser.NewReplyParser(s.reader)
	command.SendEHLO(s.writer)

	reply.GetReply(parser.Ehlo, p)

	command.SendMail(s.writer, "anurag@gmail.com")
	reply.GetReply(parser.ReplyLine, p)

	command.SendRcpt(s.writer, "anurag@gmail.com")
	reply.GetReply(parser.ReplyLine, p)

	command.SendBody(s.writer, p, "anurag@gmail.com")
	reply.GetReply(parser.ReplyLine, p)
}

func (s *Session) Begin() error {
	p := parser.NewReplyParser(s.reader)

	_, err := reply.GetReply(parser.Greeting, p)
	if err != nil {
		return err
	}
	return nil
}

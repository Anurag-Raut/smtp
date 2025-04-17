package session

import (
	"net"
	"net/http"

	"github.com/Anurag-Raut/smtp/client/dto/command"
	"github.com/Anurag-Raut/smtp/client/dto/reply"
	"github.com/Anurag-Raut/smtp/client/io/reader"
	"github.com/Anurag-Raut/smtp/client/io/writer"
	"github.com/Anurag-Raut/smtp/client/parser"
	"github.com/Anurag-Raut/smtp/logger"
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

	logger.ClientLogger.Println("EHLO SENT")
	reply.GetReply(parser.Ehlo, p)

	logger.ClientLogger.Println("EHLO GOTR REPLY")
	command.SendMail(s.writer, "anurag@gmail.com")
	reply.GetReply(parser.ReplyLine, p)

	logger.ClientLogger.Println("MAIL GOTR REPLY")

	command.SendRcpt(s.writer, "anurag@gmail.com")
	reply.GetReply(parser.ReplyLine, p)

	logger.ClientLogger.Println("RCPT GOTR REPLY")
	command.SendBody(s.writer, "anurag@gmail.com")
	reply.GetReply(parser.ReplyLine, p)
	return nil
}

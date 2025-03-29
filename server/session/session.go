package session

import (
	"bufio"

	"github.com/Anurag-Raut/smtp/server/dto/command"
	"github.com/Anurag-Raut/smtp/server/dto/reply"
	"github.com/Anurag-Raut/smtp/server/io/reader"
	"github.com/Anurag-Raut/smtp/server/parser"
)

/*
 Order of Commands
 1) EHLO
 2) MAIL
 3) RCPt



 anytime

 NOOP, HELP, EXPN, VRFY, and RSET

*/

type Session struct {
	stepIndex int
	mailStae  MailState
}

type MailState struct {
	reversePathBuffer []byte
	forwardPathBuffer []byte
	mailDataBuffer    []byte
}

func NewSession(w *bufio.Writer) *Session {

	return &Session{}
}
func (s *Session) Begin(reader *reader.Reader, writer *bufio.Writer) {
	reply.Greet(writer)
	parser := parser.NewParser(reader)
	for {
		cmd, err := command.GetCommand(parser)
		if err != nil {
			reply.HandleParseError(writer, cmd.GetCommandType(), err)
		}

	}

}

package session

import (
	"bufio"

	"github.com/Anurag-Raut/smtp/server/dto/command"
	"github.com/Anurag-Raut/smtp/server/dto/reply"
	"github.com/Anurag-Raut/smtp/server/io/reader"
	"github.com/Anurag-Raut/smtp/server/parser"
	"github.com/Anurag-Raut/smtp/server/state"
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
	mailState *state.MailState
}

func NewSession() *Session {

	return &Session{
		mailState: &state.MailState{},
	}
}
func (s *Session) Begin(reader *reader.Reader, writer *bufio.Writer) {
	reply.Greet(writer)
	parser := parser.NewParser(reader)
	for {
		cmd, err := command.GetCommand(parser)
		if err != nil {
			// reply.HandleParseError(writer, cmd.GetCommandType(), err)
		}

		responseReply := cmd.ProcessCommand(s.mailState)
		responseReply.HandleSmtpReply(writer)
	}

}

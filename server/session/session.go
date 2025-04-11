package session

import (
	"bufio"

	"github.com/Anurag-Raut/smtp/server/dto/command"
	"github.com/Anurag-Raut/smtp/server/dto/reply"
	"github.com/Anurag-Raut/smtp/server/io/reader"
	"github.com/Anurag-Raut/smtp/server/parser"
	"github.com/Anurag-Raut/smtp/server/state"
)

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
	p := parser.NewParser(reader)
	for {
		cmd, err := command.GetCommand(p)
		if err != nil {
			// reply.HandleParseError(writer, cmd.GetCommandType(), err)
		}

		replyChannel := make(chan *reply.Reply)
		go cmd.ProcessCommand(s.mailState, replyChannel)
		for {
			var responseReply *reply.Reply
			responseReply, ok := <-replyChannel
			if !ok {
				break
			}
			responseReply.HandleSmtpReply(writer)
		}
	}

}

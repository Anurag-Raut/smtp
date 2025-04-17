package session

import (
	"github.com/Anurag-Raut/smtp/logger"
	"github.com/Anurag-Raut/smtp/server/dto/command"
	"github.com/Anurag-Raut/smtp/server/dto/reply"
	"github.com/Anurag-Raut/smtp/server/io/reader"
	"github.com/Anurag-Raut/smtp/server/io/writer"
	"github.com/Anurag-Raut/smtp/server/parser"
	"github.com/Anurag-Raut/smtp/server/state"
)

type Session struct {
	mailState *state.MailState
	writer    *writer.Writer
	reader    *reader.Reader
}

func NewSession(r *reader.Reader, w *writer.Writer) *Session {

	return &Session{
		mailState: &state.MailState{},
		reader:    r,
		writer:    w,
	}
}
func (s *Session) Begin() {
	reply.Greet(s.writer)
	p := parser.NewParser(s.reader)
	for {
		logger.ServerLogger.Println(" LOOOp")
		cmd, err := command.GetCommand(p)
		if err != nil {
			logger.ServerLogger.Println("ERROR IN LOOP", err)
			return
			// reply.HandleParseError(writer, cmd.GetCommandType(), err)
		}

		logger.ServerLogger.Println("GET COMMAND TYPE", cmd.GetCommandType())

		replyChannel := make(chan reply.ReplyInterface)
		go cmd.ProcessCommand(s.mailState, replyChannel)
		for {
			var responseReply reply.ReplyInterface
			responseReply, ok := <-replyChannel
			if !ok {
				break
			}

			err := responseReply.HandleSmtpReply(s.writer)
			if err != nil {
				logger.ServerLogger.Println("EROR WHILE SENSIGN RESPONSE", err)
				break
			}
		}
	}

}

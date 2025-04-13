package command

import (
	"github.com/Anurag-Raut/smtp/logger"
	"github.com/Anurag-Raut/smtp/server/dto/reply"
	. "github.com/Anurag-Raut/smtp/server/parser"
	"github.com/Anurag-Raut/smtp/server/state"
)

type CommandInterface interface {
	GetCommandType() CommandToken
	ParseCommand() error
	ProcessCommand(mailState *state.MailState, replyChannel chan *reply.Reply)
}

type Command struct {
	commandToken CommandToken
	parser       *Parser
}

func (cmd *Command) GetCommandType() CommandToken {
	return cmd.commandToken
}

func (cmd *Command) ProcessCommand(mailState *state.MailState, replyChannel chan *reply.Reply) {
	defer close(replyChannel)
	replyChannel <- reply.NewReply(502, "Command not implemented")

}

func NewCommand(commandString string, parser *Parser) CommandInterface {
	switch commandString {
	case "EHLO":
		return &EHLO_CMD{
			Command: Command{commandToken: EHLO, parser: parser},
		}
	case "MAIL":
		return &MAIL_CMD{
			Command: Command{commandToken: MAIL},
		}

	case "RCPT":
		return &RCPT_CMD{
			Command: Command{commandToken: RCPT},
		}
	case "DATA":
		return &DATA_CMD{
			Command: Command{commandToken: EHLO},
		}
	case "NOOP":
		return &NOOP_CMD{
			Command: Command{commandToken: EHLO},
		}
	case "VRFY":
		return &VRFY_CMD{
			Command: Command{commandToken: VRFY},
		}
	case "EXPN":
		return &EXPN_CMD{
			Command: Command{commandToken: EXPN},
		}
	case "HELP":
		return &HELP_CMD{
			Command: Command{commandToken: HELP},
		}
	case "RSET":
		return &RESET_CMD{
			Command: Command{commandToken: RSET},
		}
	case "QUIT":
		return &QUIT_CMD{
			Command: Command{commandToken: QUIT},
		}
	default:
		return nil
	}

}

type EHLO_CMD struct {
	Command
	domain string
}

func (cmd *EHLO_CMD) ParseCommand() error {
	domain, err := cmd.parser.ParseEHLO()
	cmd.domain = domain
	return err
}

func (cmd *EHLO_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan *reply.Reply) {
	// send EHLO OK RSP
	defer close(replyChannel)

	err := mailState.SetMailStep(state.EHLO)
	if err != nil {
		replyChannel <- reply.NewReply(503, err.Error())
		return
	}
	mailState.ClearAll()
	replyChannel <- reply.NewReply(250)
	return
}

type MAIL_CMD struct {
	Command
	reversePath string
}

func (cmd *MAIL_CMD) ParseCommand() error {
	return nil
}
func (cmd *MAIL_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan *reply.Reply) {
	defer close(replyChannel)
	err := mailState.SetMailStep(state.MAIL)
	if err != nil {
		replyChannel <- reply.NewReply(503, err.Error())
		return
	}
	mailState.ClearAll()
	mailState.AppendReversePatahBuffer([]byte(cmd.reversePath))
	replyChannel <- reply.NewReply(250)
	return
}

type RCPT_CMD struct {
	Command
	forwardPath string
}

func (cmd *RCPT_CMD) ParseCommand() error {
	return nil
}

func (cmd *RCPT_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan *reply.Reply) {
	defer close(replyChannel)
	err := mailState.SetMailStep(state.RCPT)
	if err != nil {
		replyChannel <- reply.NewReply(503, err.Error())
		return
	}
	mailState.AppendForwardPathBuffer([]byte(cmd.forwardPath))
	replyChannel <- reply.NewReply(250)
}

type DATA_CMD struct {
	Command
	data string
}

func (cmd *DATA_CMD) ParseCommand() error {
	return nil
}

func (cmd *DATA_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan *reply.Reply) {
	defer close(replyChannel)
	err := mailState.SetMailStep(state.DATA)
	if err != nil {
		replyChannel <- reply.NewReply(503, err.Error())
		return
	}
	replyChannel <- reply.NewReply(354)
	for {
		line, err := cmd.parser.ParseDataLine()
		if err != nil {
			replyChannel <- reply.NewReply(502, err.Error())
			return
		}
		if line == "." {
			break
		}

		mailState.AppendMailDataBuffer([]byte(cmd.data))
	}
	replyChannel <- reply.NewReply(250)
}

type RESET_CMD struct {
	Command
}

func (cmd *RESET_CMD) ParseCommand() error {
	return nil
}

type VRFY_CMD struct {
	Command
	argument string
}

func (cmd *VRFY_CMD) ParseCommand() error {
	return nil
}

type EXPN_CMD struct {
	Command
	argument string
}

func (cmd *EXPN_CMD) ParseCommand() error {
	return nil
}

type HELP_CMD struct {
	Command
	argument *string
}

func (cmd *HELP_CMD) ParseCommand() error {
	return nil
}

type NOOP_CMD struct {
	Command
	argument *string
}

func (cmd *NOOP_CMD) ParseCommand() error {
	return nil
}

type QUIT_CMD struct {
	Command
}

func (cmd *QUIT_CMD) ParseCommand() error {
	return cmd.parser.ParseQuit()
}

func (cmd *QUIT_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan *reply.Reply) {
	defer close(replyChannel)
	err := mailState.SetMailStep(state.IDLE)
	if err != nil {
		replyChannel <- reply.NewReply(503)
		return
	}
	mailState.ClearAll()
	replyChannel <- reply.NewReply(221, "221 OK")
	return
}

func GetCommand(parser *Parser) (CommandInterface, error) {

	cmdToken, err := parser.ParseCommandToken()
	if err != nil {
		return nil, err
	}
	logger.ServerLogger.Println(cmdToken, ":CMD TROKEN")
	cmdObj := NewCommand(cmdToken, parser)
	err = cmdObj.ParseCommand()
	if err != nil {
		return nil, err
	}
	return cmdObj, nil
}

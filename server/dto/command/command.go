package command

import (
	"github.com/Anurag-Raut/smtp/server/dto/reply"
	. "github.com/Anurag-Raut/smtp/server/parser"
	"github.com/Anurag-Raut/smtp/server/state"
)

type CommandInterface interface {
	GetCommandType() CommandToken
	ParseCommand() error
	ProcessCommand(mailState *state.MailState) *reply.Reply
}

type Command struct {
	commandToken CommandToken
	parser       *Parser
}

func (cmd *Command) GetCommandType() CommandToken {
	return cmd.commandToken
}

func (cmd *Command) ProcessCommand(mailState *state.MailState) *reply.Reply {
	return reply.NewReply(502, "Command not implemented")
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

func (cmd *EHLO_CMD) ProcessCommand(mailState *state.MailState) *reply.Reply {
	// send EHLO OK RSP
	err := mailState.SetMailStep(state.EHLO)
	if err != nil {
		reply.NewReply(503, err.Error())
	}
	mailState.ClearAll()
	return reply.NewReply(250)
}

type MAIL_CMD struct {
	Command
	reversePath string
}

func (cmd *MAIL_CMD) ParseCommand() error {
	return nil
}
func (cmd *MAIL_CMD) ProcessCommand(mailState *state.MailState) *reply.Reply {
	err := mailState.SetMailStep(state.MAIL)
	if err != nil {
		return reply.NewReply(503, err.Error())
	}
	mailState.ClearAll()
	mailState.SetReversePathBuffer([]byte(cmd.reversePath))
	return reply.NewReply(250)
}

type RCPT_CMD struct {
	Command
	forwardPath string
}

func (cmd *RCPT_CMD) ParseCommand() error {
	return nil
}

func (cmd *RCPT_CMD) ProcessCommand(mailState *state.MailState) *reply.Reply {
	err := mailState.SetMailStep(state.RCPT)
	if err != nil {
		return reply.NewReply(503, err.Error())
	}
	mailState.SetForwardPathBuffer([]byte(cmd.forwardPath))
	return reply.NewReply(250)
}

type DATA_CMD struct {
	Command
	data string
}

func (cmd *DATA_CMD) ParseCommand() error {
	return nil
}

func (cmd *DATA_CMD) ProcessCommand(mailState *state.MailState) *reply.Reply {
	err := mailState.SetMailStep(state.DATA)
	if err != nil {
		return reply.NewReply(503, err.Error())
	}
	//TODO: make sure that we need to append or not
	mailState.SetMailDataBuffer([]byte(cmd.data))
	return reply.NewReply(250)
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
	return nil
}

func GetCommand(parser *Parser) (CommandInterface, error) {

	cmdToken, err := parser.ParseCommandToken()
	if err != nil {
		return nil, err
	}

	cmdObj := NewCommand(cmdToken, parser)
	err = cmdObj.ParseCommand()
	if err != nil {
		return nil, err
	}
	return cmdObj, nil
}

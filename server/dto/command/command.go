package command

import (
	"github.com/Anurag-Raut/smtp/server/dto/reply"
	. "github.com/Anurag-Raut/smtp/server/parser"
	"github.com/Anurag-Raut/smtp/server/state"
)

type CommandInterface interface {
	GetCommandType() CommandToken
	ParseCommand() error
	ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface)
}

type Command struct {
	commandToken CommandToken
	parser       *Parser
}

func (cmd *Command) GetCommandType() CommandToken {
	return cmd.commandToken
}

func (cmd *Command) ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface) {
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
			Command: Command{commandToken: MAIL, parser: parser},
		}

	case "RCPT":
		return &RCPT_CMD{
			Command: Command{commandToken: RCPT, parser: parser},
		}
	case "DATA":
		return &DATA_CMD{
			Command: Command{commandToken: DATA, parser: parser},
		}
	case "NOOP":
		return &NOOP_CMD{
			Command: Command{commandToken: NOOP, parser: parser},
		}
	case "VRFY":
		return &VRFY_CMD{
			Command: Command{commandToken: VRFY, parser: parser},
		}
	case "EXPN":
		return &EXPN_CMD{
			Command: Command{commandToken: EXPN, parser: parser},
		}
	case "HELP":
		return &HELP_CMD{
			Command: Command{commandToken: HELP, parser: parser},
		}
	case "RSET":
		return &RESET_CMD{
			Command: Command{commandToken: RSET, parser: parser},
		}
	case "QUIT":
		return &QUIT_CMD{
			Command: Command{commandToken: QUIT, parser: parser},
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

func (cmd *EHLO_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface) {
	defer close(replyChannel)

	err := mailState.SetMailStep(state.EHLO)
	if err != nil {
		replyChannel <- reply.NewReply(503, err.Error())
		return
	}
	mailState.ClearAll()
	replyChannel <- reply.NewEhloReply(250)
	return
}

type MAIL_CMD struct {
	Command
	reversePath string
}

func (cmd *MAIL_CMD) ParseCommand() error {
	reversepath, err := cmd.parser.ParseMail()
	if err != nil {
		return err
	}
	cmd.reversePath = reversepath
	return nil
}
func (cmd *MAIL_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface) {
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
	forwardPath, err := cmd.parser.ParseRCPT()
	if err != nil {
		return err
	}
	cmd.forwardPath = forwardPath
	return nil
}

func (cmd *RCPT_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface) {
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
	err := cmd.parser.ParseData()
	return err
}

func (cmd *DATA_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface) {
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
		mailState.AppendMailDataBuffer([]byte(line))
	}
	//TODO: store the message and then clear the state
	err = mailState.StoreBuffer()
	if err != nil {
		//TODO: check if this is correct or not
		replyChannel <- reply.NewReply(503, err.Error())
		return
	}
	mailState.ClearAll()

	replyChannel <- reply.NewReply(250)
}

type RESET_CMD struct {
	Command
}

func (cmd *RESET_CMD) ParseCommand() error {
	return cmd.parser.ParseReset()
}

func (cmd *RESET_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface) {
	defer close(replyChannel)
	mailState.ClearAll()
	mailState.SetMailStep(state.EHLO)
	replyChannel <- reply.NewReply(250, "250 OK")
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
	return cmd.parser.ParseNoop()
}

func (cmd *NOOP_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface) {
	defer close(replyChannel)
	replyChannel <- reply.NewReply(250, "250 OK")
}

type QUIT_CMD struct {
	Command
}

func (cmd *QUIT_CMD) ParseCommand() error {
	return cmd.parser.ParseQuit()
}

func (cmd *QUIT_CMD) ProcessCommand(mailState *state.MailState, replyChannel chan reply.ReplyInterface) {
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
	cmdObj := NewCommand(cmdToken, parser)
	err = cmdObj.ParseCommand()
	if err != nil {
		return nil, err
	}
	return cmdObj, nil
}

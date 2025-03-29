package command

import (
	"github.com/Anurag-Raut/smtp/server/parser"

	. "github.com/Anurag-Raut/smtp/server/parser"
)

type CommandInterface interface {
	GetCommandType() CommandToken
	ParseCommand() error
}

type Command struct {
	commandToken CommandToken
}

func (cmd *Command) GetCommandType() CommandToken {
	return cmd.commandToken
}

func NewCommand(commandString string) CommandInterface {
	switch commandString {
	case "EHLO":
		return &EHLO_CMD{
			Command: Command{commandToken: EHLO},
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
	domain, err := parser.ParseEHLO()
	return nil
}

type MAIL_CMD struct {
	Command
	reversePath string
}

func (cmd *MAIL_CMD) ParseCommand() error {
	return nil
}

type RCPT_CMD struct {
	Command
	address string
}

func (cmd *RCPT_CMD) ParseCommand() error {
	return nil
}

type DATA_CMD struct {
	Command
	data string
}

func (cmd *DATA_CMD) ParseCommand() error {
	return nil
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

func GetCommand(r *parser.Parser) (CommandInterface, error) {

	cmdToken, err := parser.ParseCommandToken()
	if err != nil {
		return nil, err
	}

	cmdObj := NewCommand(cmdToken)
	err = cmdObj.ParseCommand()
	if err != nil {
		return nil, err
	}
	return cmdObj, nil
}

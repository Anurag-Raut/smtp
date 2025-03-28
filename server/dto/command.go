package dto

import (
	"errors"
	"log"

	"github.com/Anurag-Raut/smtp/server/io/reader"
)

type CommandToken int

const (
	EHLO CommandToken = iota
	MAil
	NOT_FOUND
)

type CommandInterface interface {
	GetCommandType() CommandToken
	ParseCommand() error
}

type Command struct {
	commandToken CommandToken
}

func NewCommand(commandString string) CommandInterface {
	switch commandString {
	case "EHLO":
		return &EHLO_CMD{
			commandToken: EHLO,
		}
	default:
		return nil
	}

}

type EHLO_CMD struct {
	commandToken
	domain string
}

func (cmd *EHLO_CMD) GetCommandType() CommandToken {
	return cmd.commandToken
}
func (cmd *EHLO_CMD) ParseCommand() error {
	return nil
}

func ParseCommandToken(c CommandToken, r *reader.Reader) (CommandInterface, error) {

	cmd, err := r.GetWord(" ")
	if err != nil {
		return nil, err
	}

	cmdObj := NewCommand(cmd)
	cmdObj.ParseCommand()

}

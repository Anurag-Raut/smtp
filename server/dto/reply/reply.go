package reply

import (
	"bufio"

	"github.com/Anurag-Raut/smtp/server/dto/command"
)

var CLRF = "\r\n"

type Reply struct {
	code uint16
	text *string
}

func (r *Reply) format() []byte {
	replyString := string(r.code)
	replyString += " "
	if r.text != nil {
		replyString += (*r.text)
	}
	replyString += CLRF
	return []byte(replyString)

}

func Greet(w *bufio.Writer) error {
	text := "Anurag Server"
	rp := Reply{
		code: 220,
		text: &text,
	}
	w.Write(rp.format())
	return nil
}

func HandleParseError(w *bufio.Writer, commandToken command.CommandToken, err error) {

}

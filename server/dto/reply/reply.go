package reply

import (
	"bufio"
)

var CLRF = "\r\n"

type Reply struct {
	code uint16
	text []string
}

func (r *Reply) format() []byte {
	replyString := string(r.code)
	replyString += " "
	if r.text != nil {
		//BUG: check this out later
		replyString += (r.text[0])
	}
	replyString += CLRF
	return []byte(replyString)

}

func Greet(w *bufio.Writer) error {
	text := []string{"Anurag Server"}
	rp := Reply{
		code: 220,
		text: text,
	}
	w.Write(rp.format())
	return nil
}

func NewReply(code uint16, textlines ...string) *Reply {
	return &Reply{
		code: code,
		text: textlines,
	}
}

func (r Reply) HandleSmtpReply(w *bufio.Writer) {
}

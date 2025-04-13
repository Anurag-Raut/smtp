package reply

import (
	"bufio"
	"strconv"

	"github.com/Anurag-Raut/smtp/logger"
)

var CLRF = "\r\n"

type ReplyInterface interface {
	format() []byte
}

type Reply struct {
	code uint16
	text []string
}

type GreetingReply struct {
	Reply
	domain string
}

func (r *Reply) format() []byte {
	replyString := strconv.Itoa(int(r.code))
	if r.text != nil {

		replyString += " "
		//BUG: check this out later
		replyString += (r.text[0])
	}

	replyString += CLRF
	logger.ServerLogger.Println("reply string", replyString)
	return []byte(replyString)

}

func (r *GreetingReply) format() string {

	replyString := strconv.Itoa(int(r.code))
	replyString += " "
	replyString += r.domain
	if r.text != nil {

		replyString += " "
		//BUG: check this out later
		replyString += (r.text[0])
	}

	replyString += CLRF
	logger.ServerLogger.Println("reply greeting string", replyString)
	return (replyString)
}

func Greet(w *bufio.Writer) error {
	text := []string{"Anurag Server"}
	rp := GreetingReply{
		Reply: Reply{
			code: 220,
			text: text,
		},
		domain: "gmail.com",
	}
	_, err := w.WriteString(rp.format())
	if err != nil {
		logger.ServerLogger.Println(err, "ERROR")
		return err
	}
	return w.Flush()
}

func NewReply(code uint16, textlines ...string) *Reply {
	return &Reply{
		code: code,
		text: textlines,
	}
}

func (r Reply) HandleSmtpReply(w *bufio.Writer) {
}

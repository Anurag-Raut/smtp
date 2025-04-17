package reply

import (
	"github.com/Anurag-Raut/smtp/logger"
	"github.com/Anurag-Raut/smtp/server/io/writer"
	"strconv"
)

var CLRF = "\r\n"

type ReplyInterface interface {
	format() string
	HandleSmtpReply(w *writer.Writer) error
}

type Reply struct {
	code uint16
	text []string
}

type GreetingReply struct {
	Reply
	domain string
}

type EhloReply struct {
	Reply
	domain string
}

func (r *Reply) format() string {
	replyString := strconv.Itoa(int(r.code))
	if r.text != nil {

		replyString += " "
		//BUG: check this out later
		replyString += (r.text[0])
	}

	replyString += CLRF
	logger.ServerLogger.Println("reply string", replyString)
	return replyString

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

func (r *EhloReply) format() string {

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
func Greet(w *writer.Writer) error {
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

func NewReply(code uint16, textlines ...string) ReplyInterface {
	return &Reply{
		code: code,
		text: textlines,
	}
}

func NewEhloReply(code uint16, textlines ...string) ReplyInterface {
	return &EhloReply{
		Reply: Reply{

			code: code,
			text: textlines,
		},
		domain: "gmail.com",
	}
}

func (r *Reply) HandleSmtpReply(w *writer.Writer) error {
	_, err := w.WriteString(r.format())
	return err
}

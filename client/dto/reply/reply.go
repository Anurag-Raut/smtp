package reply

import (
	"errors"

	"github.com/Anurag-Raut/smtp/client/parser"
	"github.com/Anurag-Raut/smtp/logger"
)

type Reply struct {
	code        int
	textStrings []string

	parser *parser.ReplyParser
}

type ReplyInterface interface {
	ParseReply() error
	Execute() error
}

type GreetingReply struct {
	serverIdentifier string
	Reply
}
type EhloReply struct {
	domain string
	Reply
}

func (r *GreetingReply) ParseReply() error {
	identifier, textStrings, err := r.parser.ParseGreeting()
	if err != nil {
		return err
	}
	r.code = 220
	r.textStrings = textStrings
	r.serverIdentifier = identifier
	return nil
}

func (r *EhloReply) ParseReply() error {
	domain, textStrings, err := r.parser.ParseGreeting()
	if err != nil {
		return err
	}
	r.code = 220
	r.textStrings = textStrings
	r.domain = domain
	return nil
}

func (r *Reply) ParseReply() error {
	replyCode, textStrings, err := r.parser.ParseReplyLine()
	if err != nil {
		return err
	}
	r.code = replyCode
	r.textStrings = textStrings
	return nil
}

func (r *Reply) Execute() error {
	return nil
}

func (r *GreetingReply) Execute() error {
	return nil
}
func GetReply(token parser.ReplyToken, p *parser.ReplyParser) (reply ReplyInterface, err error) {

	switch token {
	case parser.ReplyLine:
		reply = &Reply{
			parser: p,
		}
		break
	case parser.Greeting:
		reply = &GreetingReply{
			Reply: Reply{parser: p},
		}
	case parser.Ehlo:
		reply = &EhloReply{
			Reply: Reply{parser: p},
		}
	default:
		{
			return nil, errors.New("Could not find the Reply")
		}

	}
	logger.ClientLogger.Println("PRINTING reply tyoe", reply)
	err = reply.ParseReply()
	if err != nil {
		return nil, err
	}
	return reply, nil
}

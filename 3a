package reply

import "github.com/Anurag-Raut/smtp/client/parser"

type Reply struct {
	code   int
	parser *parser.ReplyParser
}

type ReplyInterface interface {
	ParseReply() error
}

type GreetingReply struct {
	Reply
	textStrings      []string
	serverIdentifier string
}

func (r *GreetingReply) ParseReply() error {
	identifier, textStrings, err := r.parser.ParseGreeting()
	if err != nil {
		return err
	}
	r.textStrings = textStrings
	r.serverIdentifier = identifier
	return nil
}

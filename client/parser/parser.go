package parser

import (
	"bufio"
	"errors"
	"strings"

	"github.com/Anurag-Raut/smtp/client/io/reader"
)

func GetDomainFromEmail(email string) (domain string, err error) {
	splits := strings.Split(email, "@")
	if len(splits) != 2 {
		return "", errors.New("Invalid Email , domain couldn't be parsed")
	}
	return splits[1], nil
}

type ReplyParser struct {
	reader reader.Reader
}

func NewReplyParser(r *bufio.Reader) *ReplyParser {
	return &ReplyParser{reader: *reader.NewReader(r)}
}

func (p *ReplyParser) parseGreeting() error {
	statusCode, err := p.reader.ReadStringOfLen(3)
	if err != nil {
		return err
	}

	return nil
}

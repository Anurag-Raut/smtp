package parser

import (
	"errors"
	"fmt"
	"text/template/parse"
	"unicode"

	"github.com/Anurag-Raut/smtp/server/io/reader"
)

type Parser struct {
	reader *reader.Reader
}

func NewParser(reader *reader.Reader) *Parser {
	return &Parser{
		reader: reader,
	}
}

func (p *Parser) ParseCommandToken() (string, error) {
	return p.reader.ReadStringOfLen(4)
}

type TokenNotFound struct {
	token TokenType
}

func (t TokenNotFound) Error() string {
	return fmt.Sprintf("Token not found: %d", t.token)
}

func (p *Parser) ExpectMultiple(tokens ...TokenType) (string, error) {
	for _, token := range tokens {
		value, err := p.Expect(token)
		if err == nil {
			return value, nil
		}
	}
	return "", TokenNotFound{}

}
func (p *Parser) Expect(token TokenType) (string, error) {
	switch token {
	case SPACE:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}
			if string(bytes) != " " {
				return "", TokenNotFound{token: token}
			}
			_, err = p.reader.ReadByte()
			if err != nil {
				return "", err
			}
			return string(bytes), nil

		}
	case LEFT_ANGLE_BRAC:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}
			if string(bytes) != "<" {
				return "", TokenNotFound{token: token}
			}
			_, err = p.reader.ReadByte()
			if err != nil {
				return "", err
			}
			return string(bytes), nil

		}
	case RIGHT_ANGLE_BRAC:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}
			if string(bytes) != ">" {
				return "", TokenNotFound{token: token}
			}
			_, err = p.reader.ReadByte()
			if err != nil {
				return "", err
			}
			return string(bytes), nil

		}
	case ALPHA:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}

			if !unicode.IsLetter(rune(bytes[0])) {
				return "", TokenNotFound{token: token}
			}
			_, err = p.reader.ReadByte()
			if err != nil {
				return "", err
			}
			return string(bytes), nil

		}

	}

	return "", TokenNotFound{}
}

func (p *Parser) ParseEHLO() (domain string, err error) {
	_, err = p.Expect(SPACE)
	if err != nil {
		return "", err
	}
	domain, err = p.parseDomain()
	return domain, err
}

func (p *Parser) parseDomain() (string, error) {
	_, err := p.Expect(LEFT_ANGLE_BRAC)
	if err != nil {
		return "", err
	}
	subDomain, err := p.parseSubDomain()
	for {
		_, err := p.Expect(DOT)
		if err != nil {
			if (errors.Is(err, TokenNotFound{})) {
				break
			} else {
				return "", err
			}
		}
		subDomain += "."
		newSubDomain, err := p.parseSubDomain()
		if err != nil {
			return "", err
		}
		subDomain += newSubDomain
	}
	return subDomain, nil
}

func (p *Parser) parseSubDomain() (string, error) {
	firstVal, err := p.ExpectMultiple(ALPHA, DIGIT)
	if err != nil {
		return "", err
	}
	middleVal := ""
	for {
		ch, err := p.ExpectMultiple(ALPHA, DIGIT)
		if err != nil {
			if (errors.Is(err, TokenNotFound{})) {
				break
			} else {
				return "", err
			}
		}
		middleVal += ch
	}

	if len(middleVal) > 0 {
		err := p.reader.UnreadByte()
		if err != nil {
			return "", err
		}
		_, err = p.ExpectMultiple(ALPHA, DIGIT)

		if err != nil {
			return firstVal + middleVal, err
		}

	}
	return firstVal + middleVal, nil

}

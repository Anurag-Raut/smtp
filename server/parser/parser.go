package parser

import (
	"errors"
	"fmt"
	"strings"
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
	case AT:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}
			if string(bytes) != "@" {
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

func (p *Parser) ParseEHLO() (string, error) {
	_, err := p.Expect(SPACE)
	if err != nil {
		return "", err
	}
	domain, err := p.parseDomain()
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

func (p *Parser) ParseMail() (string, error) {
	_, err := p.Expect(SPACE)
	if err != nil {
		return "", err
	}
	fromString, err := p.reader.ReadStringOfLen(4)
	if err != nil {
		return "", err
	}
	if strings.ToLower(fromString) != "from" {
		return "", TokenNotFound{token: KEYWORD}
	}
	reversePath, err := p.parseReversePath()
	return reversePath, nil
}

func (p *Parser) parseReversePath() (string, error) {
	path, err := p.parsePath()
	if err == nil {
		return path, err
	}

	_, err = p.Expect(LEFT_ANGLE_BRAC)
	if err != nil {
		return "", err
	}
	_, err = p.Expect(LEFT_ANGLE_BRAC)
	if err != nil {
		return "", err
	}

	return "", nil
}

func (p *Parser) parsePath() (string, error) {
	_, err := p.Expect(LEFT_ANGLE_BRAC)
	if err != nil {
		return "", err
	}

	err = p.parseAD1() //ignore source routes

	if err != nil {
		return "", err
	}
	mailbox, err := p.parseMailBox()
	_, err = p.Expect(RIGHT_ANGLE_BRAC)
	if err != nil {
		return "", err
	}

	return mailbox, nil

}

func (p *Parser) parseAD1() error {

	_, err := p.Expect(AT)
	if err != nil {
		return nil
	}

	_, err = p.parseDomain()
	if err != nil {
		return err
	}
	for {
		_, err := p.Expect(AT)
		if err != nil {
			break
		}

		_, err = p.parseDomain()
		if err != nil {
			return err
		}
	}

	_, err = p.Expect(COLON)
	if err != nil {
		return err
	}

	return nil

}

func (p *Parser) parseMailBox() (string, error) {
	localPart, err := p.parseLocalPart()
	if err != nil {
		return "", err
	}
	_, err = p.Expect(AT)
	if err != nil {
		return "", err
	}
	domain, err := p.parseDomain()
	if err == nil {

		return localPart + "@" + domain, nil
	}
	addressLiteral, err := p.parseAddressLiteral()
	if err != nil {
		return "", err
	}

	return localPart + "@" + addressLiteral, nil

}

func (p *Parser) parseLocalPart() (string, error) {
	dotString, err := p.parseDotString()
	if err == nil {
		return dotString, nil
	}
	quotedString, err := p.parseQuotedString()
	if err != nil {
		return "", err
	}

	return quotedString, nil
}

func (p *Parser) parseDotString() (string, error) {
	atom, err := p.parseAtom()
	if err != nil {
		return "", err
	}
	value := atom
	for {
		_, err := p.Expect(AT)
		if err != nil {
			break
		}

		atom, err := p.parseAtom()
		if err != nil {
			return "", err
		}
		value += atom + "." + atom
	}

	return value, nil
}

func (p *Parser) parseAtom() (string, error) {

	atom := ""
	ch, err := p.Expect(ATEXT)
	if err != nil {
		return "", err
	}
	atom += ch
	for {
		ch, err := p.ExpectMultiple(ATEXT)
		if err != nil {
			break
		}
		atom += ch
	}
	return atom, nil
}
func (p *Parser) parseQuotedString() (string, error) {
	_, err := p.Expect(DQUOTE)
	if err != nil {
		return "", err
	}
	value := ""

	for {
		ch, err := p.Expect(QTEXTSMTP)
		if err != nil {
			value = ""
			break
		}
		value += ch
	}

	if value == "" {
		for {
			ch, err := p.Expect(QPAIRSMTP)
			if err != nil {
				value = ""
				break
			}
			value += ch
		}

	}

	_, err = p.Expect(DQUOTE)
	if err != nil {
		return "", err
	}

	return string('"') + value + string('"'), nil
}

func (p *Parser) parseAddressLiteral() (string, error) {
	return "", nil
}

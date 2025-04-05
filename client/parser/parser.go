package parser

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"unicode"

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

func (p *ReplyParser) ParseGreeting() (identifier string, textStrings []string, err error) {
	statusCodeString, err := p.expect(CODE)
	if err != nil {
		return identifier, textStrings, err
	}
	statusCode, err := strconv.Atoi(statusCodeString)
	if err != nil {
		return identifier, textStrings, errors.New("Error while converting the status code into int")
	}
	if statusCode != 220 {
		return identifier, textStrings, errors.New("Expected status Code as 220")
	}

	_, err = p.expect(SPACE)
	if err == nil {
		identifier, textString, err := p.parseSingleLine()
		if err != nil {
			return identifier, textString, err
		}
		return identifier, textString, nil

	} else if (errors.Is(err, TokenNotFound{})) {

	}

	return identifier, textStrings, nil
}

func (p *ReplyParser) parseSingleLine() (identifier string, textStrings []string, err error) {
	identifier, err = p.parseDomain()
	if err != nil {
		identifier, err = p.parseAddressLiteral()
		if err != nil {
			return identifier, textStrings, err
		}

	}

	_, err = p.expect(SPACE)
	textString := ""
	if err == nil {
		textString, err = p.parserTextString()
		if err != nil {
			return identifier, textStrings, err
		}

	}
	_, err = p.expect(CRLF)

	if err != nil {
		return identifier, textStrings, err
	}

	return identifier, []string{textString}, nil

}

func (p *ReplyParser) parseMultiLineTextString() (identifier string, textStrings []string, err error) {

	identifier, err = p.parseDomain()

	if err != nil {
		identifier, err = p.parseAddressLiteral()
		if err != nil {
			return identifier, textStrings, err
		}

	}

	_, err = p.expect(SPACE)
	if err == nil {
		textString, err := p.parserTextString()
		if err != nil {
			return identifier, textStrings, err
		}
		textStrings = append(textStrings, textString)
	}
	for {
		codeString, err := p.expect(CODE)
		if err != nil {
			return identifier, textStrings, err
		}
		code, err := strconv.Atoi(codeString)
		if err != nil {
			return identifier, textStrings, err
		}
		if code != 220 {
			return identifier, textStrings, errors.New("Could not parse Code")
		}

		_, err = p.expect(HYPHEN)
		if err == nil {
			textStr, err := p.parserTextString()
			if err != nil {
				return identifier, textStrings, err
			}
			textStrings = append(textStrings, textStr)
			_, err = p.expect(CRLF)
			if err != nil {
				return identifier, textStrings, err
			}
		} else if (errors.Is(err, TokenNotFound{})) {
			break
		} else {
			return identifier, textStrings, err
		}
	}
	_, err = p.expect(SPACE)

	if err == nil {
		textString, err := p.parserTextString()
		if err != nil {
			return identifier, textStrings, err
		}
		textStrings = append(textStrings, textString)
	}
	_, err = p.expect(CRLF)
	if err != nil {
		return identifier, textStrings, err

	}

	return identifier, textStrings, nil
}
func (p *ReplyParser) parserTextString() (string, error) {

	var textString string
	for {
		ch, err := p.expectMultiple(ALPHA, SPACE, HT)
		if err == nil {
			textString += ch
		} else if (errors.Is(err, TokenNotFound{})) {
			break
		} else {
			return textString, err
		}
	}

	return textString, nil
}

func (p *ReplyParser) parseDomain() (string, error) {
	_, err := p.expect(LEFT_ANGLE_BRAC)
	if err != nil {
		return "", err
	}
	subDomain, err := p.parseSubDomain()
	for {
		_, err := p.expect(DOT)
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

func (p *ReplyParser) parseSubDomain() (string, error) {
	firstVal, err := p.expectMultiple(ALPHA, DIGIT)
	if err != nil {
		return "", err
	}
	middleVal := ""
	for {
		ch, err := p.expectMultiple(ALPHA, DIGIT)
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
		_, err = p.expectMultiple(ALPHA, DIGIT)

		if err != nil {
			return firstVal + middleVal, err
		}

	}
	return firstVal + middleVal, nil

}

func (p *ReplyParser) parseAddressLiteral() (string, error) {
	_, err := p.expect(LEFT_SQUARE_BRAC)
	if err != nil {
		return "", err
	}

	ip4addres, err := p.parseIPV4_AddressLiteral()
	if err != nil {
		return "", err
	}
	_, err = p.expect(RIGHT_SQUARE_BRAC)
	if err != nil {
		return "", err
	}
	return ip4addres, nil
}

func (p *ReplyParser) parseIPV4_AddressLiteral() (string, error) {
	ipv4_address := ""
	for i := 0; i < 3; i++ {
		val, err := p.expect(DIGIT)
		if err != nil {
			return "", err
		}
		ipv4_address += val
	}

	for j := 0; j < 3; j++ {
		dotString, err := p.expect(DOT)
		if err != nil {
			return "", err
		}
		ipv4_address += dotString
		for i := 0; i < 3; i++ {
			val, err := p.expect(DIGIT)
			if err != nil {
				return "", err
			}
			ipv4_address += val
		}
	}

	return ipv4_address, nil

}

func (p *ReplyParser) expectMultiple(tokens ...TokenType) (string, error) {
	for _, token := range tokens {
		value, err := p.expect(token)
		if err == nil {
			return value, nil
		}
	}
	return "", TokenNotFound{}

}

func (p *ReplyParser) expect(token TokenType) (string, error) {
	switch token {

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

	case DIGIT:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}

			if !unicode.IsDigit(rune(bytes[0])) {
				return "", TokenNotFound{token: token}
			}
			_, err = p.reader.ReadByte()
			if err != nil {
				return "", err
			}
			return string(bytes), nil

		}
	case CODE:
		{

			bytes, err := p.reader.Peek(3)
			if err != nil {
				return "", nil
			}
			code := string(bytes)
			first := code[0]
			second := code[1]
			third := code[2]

			if first < 0x32 || first > 0x35 {
				return "", errors.New("Expected first value of code to be between 0x32 and 0x35")
			}

			if second < 0x30 || second > 0x35 {
				return "", errors.New("Expected second value of code to be between 0x30 and 0x35")
			}

			if third < 0x30 || third > 0x39 {
				return "", errors.New("Expected third value of code to be between 0x30 and 0x39")
			}
			byteArray := make([]byte, 3)
			_, err = p.reader.Read(byteArray)
			if err != nil {
				return "", err
			}
			return code, nil
		}
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
	case CRLF:
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

	default:
		{
			return "", errors.New("No Matching token found")
		}
	}

}

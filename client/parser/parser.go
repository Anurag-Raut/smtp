package parser

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/Anurag-Raut/smtp/client/io/reader"
	"github.com/Anurag-Raut/smtp/logger"
)

func GetDomainFromEmail(email string) (domain string, err error) {
	splits := strings.Split(email, "@")
	if len(splits) != 2 {
		return "", errors.New("Invalid Email , domain couldn't be parsed")
	}
	return splits[1], nil
}

type ReplyParser struct {
	reader *reader.Reader
}

func NewReplyParser(r *reader.Reader) *ReplyParser {
	return &ReplyParser{reader: r}
}

func (p *ReplyParser) ParseGreeting() (identifier string, textStrings []string, err error) {
	statusCodeString, err := p.expect(CODE)
	logger.ClientLogger.Println(statusCodeString)
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
		return identifier, textString, err

	} else {
		_, err := p.expect(HYPHEN)
		if err != nil {
			return identifier, textStrings, err
		}
		return p.parseMultiLineTextString()

	}

}
func (p *ReplyParser) ParseReplyLine() (replyCode int, textStrings []string, err error) {

	codeString, err := p.expect(CODE)
	if err != nil {
		return replyCode, textStrings, err
	}

	replyCode, err = strconv.Atoi(codeString)
	if err != nil {
		return replyCode, textStrings, err
	}
	for {
		codeString, err := p.expect(CODE)
		if err != nil {
			return replyCode, textStrings, err
		}

		rCode, err := strconv.Atoi(codeString)
		if err != nil {
			return replyCode, textStrings, err
		}

		if rCode != replyCode {
			return replyCode, textStrings, errors.New("Expected same reply code as first line b ut got: " + codeString)
		}

		_, err = p.expect(HYPHEN)
		if err == nil {
			txtString, err := p.parserTextString()
			if err != nil {
				return replyCode, textStrings, nil
			}
			textStrings = append(textStrings, txtString)
			_, err = p.expect(CRLF)
			if err != nil {
				return replyCode, textStrings, nil
			}
		} else {

			_, err := p.expect(SPACE)
			if err == nil {
				txtString, err := p.parserTextString()
				if err != nil {
					return replyCode, textStrings, err
				}
				textStrings = append(textStrings, txtString)

			}
			_, err = p.expect(CRLF)
			if err != nil {
				return replyCode, textStrings, err
			}
			return replyCode, textStrings, nil
		}

	}

}

func (p *ReplyParser) parseSingleLine() (identifier string, textStrings []string, err error) {
	identifier, err = p.parseDomain()
	logger.ClientLogger.Println("ERROR IN Parsing domain", err)
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
		logger.ClientLogger.Println("PARSE TEXT STRING", textString)
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
	subDomain, err := p.parseSubDomain()
	if err != nil {
		return "", err
	}
	for {
		_, err := p.expect(DOT)
		if err != nil {
			if (errors.Is(err, TokenNotFound{token: DOT})) {
				break
			} else {
				return "", err
			}

		}
		logger.ClientLogger.Println("we already here: ", subDomain)
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
		ch, err := p.expectMultiple(ALPHA, DIGIT, HYPHEN)
		if err != nil {
			if (errors.Is(err, TokenNotFound{})) {
				break
			} else {
				return "", err
			}
		}
		middleVal += ch
	}
	// byh, _ := p.reader.Peek(1)
	// logger.ClientLogger.Println("PRITING BYTE AFTER in subdomain, ", string(byh))

	// if len(middleVal) > 0 {
	//
	// 	_, err = p.expectMultiple(ALPHA, DIGIT)
	//
	// 	if err != nil {
	// 		return firstVal + middleVal, err
	// 	}
	//
	// }
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

			logger.ClientLogger.Println("EXPECTING CODE")
			bytes, err := p.reader.Peek(3)
			logger.ClientLogger.Println("PEEKIED", bytes)
			logger.ClientLogger.Println(string(bytes), "BYTES")

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
	case DOT:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}
			if string(bytes) != "." {
				return "", TokenNotFound{token: token}
			}
			_, err = p.reader.ReadByte()
			if err != nil {
				return "", err
			}
			return string(bytes), nil
		}

	case HYPHEN:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}
			if string(bytes) != "-" {
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
			bytes, err := p.reader.Peek(2)
			if err != nil {
				return "", err
			}
			if string(bytes) != "\r\n" {
				return "", TokenNotFound{token: token}
			}
			crlfBytes := make([]byte, 2)
			_, err = p.reader.Read(crlfBytes)
			if err != nil {
				return "", err
			}
			return string(crlfBytes), nil
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
	case LEFT_SQUARE_BRAC:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}

			if string(bytes) != "[" {
				return "", TokenNotFound{token: token}
			}
			_, err = p.reader.ReadByte()
			if err != nil {
				return "", err
			}
			return string(bytes), nil

		}
	case RIGHT_SQUARE_BRAC:
		{
			bytes, err := p.reader.Peek(1)
			if err != nil {
				return "", err
			}
			if string(bytes) != "]" {
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
			return "", errors.New("No Matching token found " + strconv.Itoa(int(token)))
		}
	}

}

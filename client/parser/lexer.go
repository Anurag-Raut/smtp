package parser

import "fmt"

type TokenType int

const (
	DIGIT TokenType = iota
	ALPHA
	CODE
	SPACE
	CRLF
	DOT
	LEFT_ANGLE_BRAC
	RIGHT_ANGLE_BRAC
	LEFT_SQUARE_BRAC
	RIGHT_SQUARE_BRAC
	HT
	HYPHEN
)

type TokenNotFound struct {
	token TokenType
}

func (t TokenNotFound) Error() string {
	return fmt.Sprintf("Token Not found:  %d ", t.token)
}

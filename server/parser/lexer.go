package parser

type TokenType int

const (
	// COMMANDS

	TEXT TokenType = iota
	CRLF
	SPACE
	LEFT_ANGLE_BRAC
	RIGHT_ANGLE_BRAC
	COLON
	ALPHA
	DIGIT
	HYPHEN
	DOT
	KEYWORD
	AT
	ATEXT
	QTEXTSMTP
	QPAIRSMTP
	DQUOTE
)

type CommandToken TokenType

const (
	EHLO CommandToken = iota
	HELO
	MAIL
	RCPT
	QUIT
	EXPN
	VRFY
	NOOP
	DATA
	RSET
	HELP
	NOT_FOUND
)

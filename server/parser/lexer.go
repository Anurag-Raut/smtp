package parser

type TokenType int

const (
	// COMMANDS
	EHLO TokenType = iota
	HELO
	CLRF
	SP
	DOMAIN
	MAIl
	FROM
	COLON
	D_QUOTE
	// Reply
	CODE
	REPLYTEXT

	TEXT
)

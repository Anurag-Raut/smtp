package command

import (
	"fmt"

	"github.com/Anurag-Raut/smtp/client/config"
	"github.com/Anurag-Raut/smtp/client/dto/reply"
	"github.com/Anurag-Raut/smtp/client/io/writer"
	"github.com/Anurag-Raut/smtp/client/parser"
)

func SendEHLO(w *writer.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("EHLO %s\r\n", config.ClientDomain))
	return err
}

func SendMail(w *writer.Writer, reversePath string) error {
	_, err := w.WriteString(fmt.Sprintf("MAIL FROM:<%s>\r\n", reversePath))
	return err
}

func SendRcpt(w *writer.Writer, forwardPath string) error {
	_, err := w.WriteString(fmt.Sprintf("RCPT TO:<%s>\r\n", forwardPath))
	return err
}

func SendBody(w *writer.Writer, p *parser.ReplyParser, body string) error {

	_, err := w.WriteString(fmt.Sprintf("DATA\r\n"))
	if err != nil {
		return err
	}

	_, err = reply.GetReply(parser.ReplyLine, p)
	if err != nil {
		return err
	}

	_, err = w.WriteString(fmt.Sprintf("%s\r\n", body))

	if err != nil {
		return err
	}

	_, err = w.WriteString(fmt.Sprintf(".\r\n"))

	if err != nil {
		return err
	}
	return nil
}

func SendReset(w *writer.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("RSET\r\n"))

	return err
}

func SendVerify(w *writer.Writer, identifier string) error {
	_, err := w.WriteString(fmt.Sprintf("VRFY %s\r\f", identifier))
	return err
}

func SendExpand(w *writer.Writer, mailingList string) error {
	_, err := w.WriteString(fmt.Sprintf("EXPN %s\r\n", mailingList))
	return err
}
func SendHelp(w *writer.Writer, argument *string) error {
	if argument != nil {
		_, err := w.WriteString(fmt.Sprintf("HELP %s\r\n", *argument))
		return err
	} else {
		_, err := w.WriteString(fmt.Sprintf("HELP\r\n"))
		return err
	}

}

func SendNoop(w *writer.Writer) error {

	_, err := w.WriteString(fmt.Sprintf("NOOP\r\n"))
	return err
}

func SendQuit(w *writer.Writer) error {
	_, err := w.WriteString(fmt.Sprintf("QUIT\r\f"))
	return err
}

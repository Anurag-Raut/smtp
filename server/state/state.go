package state

import (
	"errors"

	"github.com/Anurag-Raut/smtp/server/store"
)

type MAIL_STEP int

const (
	IDLE MAIL_STEP = iota
	EHLO
	MAIL
	RCPT
	DATA
)

type MailState struct {
	reversePathBuffer []byte
	forwardPathBuffer []byte
	mailDataBuffer    []byte

	stepIndex MAIL_STEP
}

func (mailState *MailState) ClearReversePathBuffer() {
	mailState.reversePathBuffer = make([]byte, 0)
}

func (mailState *MailState) ClearForwardPathBuffer() {
	mailState.forwardPathBuffer = make([]byte, 0)
}

func (mailState *MailState) ClearMailDataBuffer() {
	mailState.mailDataBuffer = make([]byte, 0)
}

func (mailState *MailState) ClearAll() {
	mailState.ClearMailDataBuffer()
	mailState.ClearReversePathBuffer()
	mailState.ClearForwardPathBuffer()
}

func (mailState *MailState) AppendReversePatahBuffer(data []byte) {
	mailState.reversePathBuffer = data
}

func (mailState *MailState) AppendForwardPathBuffer(data []byte) {
	mailState.forwardPathBuffer = data
}

func (mailState *MailState) AppendMailDataBuffer(data []byte) {
	mailState.mailDataBuffer = data
}

func (mailState *MailState) SetMailStep(step MAIL_STEP) error {
	if step == EHLO {

	} else if mailState.stepIndex >= step {
		msg := "Bad sequence of commands"
		return errors.New(msg)
	}
	mailState.stepIndex = step
	return nil
}

func (mailState *MailState) StoreBuffer() error {
	if mailState.stepIndex != DATA {
		return errors.New("ERROR, invalid store command , step not set to store")
	}
	err := store.StoreEmail(string(mailState.reversePathBuffer), string(mailState.forwardPathBuffer), string(mailState.mailDataBuffer))
	return err

}

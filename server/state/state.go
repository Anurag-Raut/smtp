package state

import (
	"errors"
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

	stepIndex int
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

func (mailState *MailState) SetReversePathBuffer(data []byte) {
	mailState.reversePathBuffer = data
}

func (mailState *MailState) SetForwardPathBuffer(data []byte) {
	mailState.forwardPathBuffer = data
}

func (mailState *MailState) SetMailDataBuffer(data []byte) {
	mailState.mailDataBuffer = data
}

func (mailState *MailState) SetMailStep(step MAIL_STEP) error {
	if step == EHLO {

	} else if mailState.stepIndex >= int(step) {
		msg := "Bad sequence of commands"
		return errors.New(msg)
	}
	mailState.stepIndex = int(step)
	return nil
}

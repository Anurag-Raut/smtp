package writer

import (
	"bufio"
	"net"
)

type Writer struct {
	*bufio.Writer
}

func NewWriter(conn net.Conn) *Writer {
	return &Writer{Writer: bufio.NewWriter(conn)}
}

func (w *Writer) WriteString(s string) (n int, err error) {
	n, err = w.Writer.WriteString(s)
	if err != nil {
		return n, err
	}
	err = w.Flush()

	return n, err
}

package reader

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type Reader struct {
	*bufio.Reader
}

func NewReader(conn net.Conn) *Reader {

	bufIOReader := bufio.NewReader(conn)
	reader := Reader{bufIOReader}
	return &reader
}

func (r *Reader) GetLine(delim string) (string, error) {
	var word string = ""
	var delimIndex int = 0
	for {

		ch_bytes, err := r.Peek(1)
		if err != nil {
			return "", err
		}

		ch := string(ch_bytes)

		if string(delim[delimIndex]) == ch {
			for {
				potential_delim_bytes, err := r.Peek(delimIndex + 1)
				if err != nil {
					return word, errors.New("DID not found delim" + err.Error())
				}
				next_ch := string(potential_delim_bytes[delimIndex])
				if string(delim[delimIndex]) == next_ch {
					delimIndex++

				} else {
					delimIndex = 0
					break
				}

				if delimIndex == len(delim) {
					bytesToRead := make([]byte, len(delim))
					_, err = r.Read(bytesToRead)
					return word, err
				}
			}
		} else {
			delimIndex = 0
			_, err := r.ReadByte()
			if err != nil {
				return "", err
			}
			word += ch
		}

	}
	return word, nil
}

func (r Reader) ReadStringOfLen(n int) (string, error) {
	var cmdBytes []byte = make([]byte, n)

	readLen, err := r.Read(cmdBytes)
	if err != nil {
		return "", err
	}
	if readLen != n {
		return "", errors.New(fmt.Sprintf("Could not get string of %d bytes", n))

	}
	return string(cmdBytes), nil
}

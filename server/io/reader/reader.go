package reader

import (
	"bufio"
	"errors"
)

type Reader struct {
	*bufio.Reader
	index int
}

func (r *Reader) GetWord(delim string) (string, error) {
	var word string = ""
	var delimIndex int = 0
	for {
		ch_bytes, err := r.Peek(1)
		if err != nil {
			return word, err
		}
		ch := string(ch_bytes)
		if string(delim[delimIndex]) == ch {
			for {
				potential_delim_bytes, err := r.Peek(delimIndex + 1)
				if err != nil {
					return word, errors.New("DID not found delim" + err.Error())
				}
				next_ch := string(potential_delim_bytes[delimIndex])
				if delimIndex == len(delim) {
					return word, nil
				} else if string(delim[delimIndex]) == next_ch {
					delimIndex++
				} else {
					delimIndex = 0
					break
				}
			}
		} else {
			delimIndex = 0
		}

	}
	return word, nil
}

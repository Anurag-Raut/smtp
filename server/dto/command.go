package dto

import "bufio"

type Command int

const (
	EHLO Command = iota
	MAil
)

// func readLine(reader bufio.Reader) (string, error) {
// 	var s string = ""
// 	for {
// 		r, _, err := reader.ReadRune()
// 		if err != nil {
// 			return "", err
// 		}
//
// 		if r == '\r' {
// 			s += string(r)
// 			nextRn, _, err := reader.ReadRune()
// 			if err != nil {
// 				return "", err
// 			}
// 			s += string(nextRn)
// 			if nextRn == '\n' {
// 				return s, nil
// 			}
// 		} else {
// 			s += string(r)
// 		}
// 	}
//
//   return s,nil
//
// }

func 

func ParseCommand(c Command, r *bufio.Reader, w *bufio.Writer) {
  line,err:=
}

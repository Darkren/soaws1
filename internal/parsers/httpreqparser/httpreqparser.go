package httpreqparser

import (
	"fmt"
	"soawstest/internal/parsers/utility"
	"strconv"
	"strings"
)

type HttpRequestParser struct {
	payload []byte
}

func New() *HttpRequestParser {
	return &HttpRequestParser{[]byte{}}
}

func (hrp *HttpRequestParser) Next() []byte {
	current := 0
	contentLength := 0
	var line string
	var err error
	for {
		prev := current
		current = utility.Seek(&hrp.payload, current, 0x0A)

		//	package is not complete yet
		if current == len(hrp.payload) {
			return nil
		}

		line = string(hrp.payload[prev : current-1])
		if strings.HasPrefix(line, "Content-Length") {
			contentLength, err = strconv.Atoi(strings.Trim(strings.Split(line, ":")[1], " "))
			if err != nil {
				fmt.Println("Content-Length header mailformed")

				contentLength = 0
			}
		}

		current++
		//	package is not complete yet
		if (current + 1) >= len(hrp.payload) {
			return nil
		}

		if hrp.payload[current] == 0x0D {
			var req []byte
			if contentLength == 0 {
				req = hrp.payload[0 : current+2]
			} else {

				current += contentLength
				//	package is not complete yet
				if current >= len(hrp.payload) {
					return nil
				}

				req = hrp.payload[0 : current+2]
			}

			hrp.payload = hrp.payload[current+2:]

			return req
		}
	}
}

// Append is used to append new portion of data to the local payload buffer
func (hrp *HttpRequestParser) Append(data *[]byte) {
	hrp.payload = append(hrp.payload, (*data)...)
}

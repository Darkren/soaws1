package requests

import (
	"fmt"
	"soawstest/internal/parsers/utility"
	"strconv"
	"strings"
)

type HttpRequest struct {
	ContentLength int
	Method        string
	Body          string
}

func (hr *HttpRequest) New(data *[]byte) *HttpRequest {
	current := utility.Seek(data, 0, 0x20)

	hr.Method = string((*data)[0:current])

	var line string
	var err error
	isProcessingFinished := false
	for !isProcessingFinished {
		prev := current

		current = utility.Seek(data, current, 0x0D)

		line = string((*data)[prev:current])
		if strings.HasPrefix(line, "Content-Length") {
			hr.ContentLength, err = strconv.Atoi(
				strings.Trim(strings.Split(line, ":")[1], " "))

			if err != nil {
				fmt.Printf("Content-Length header mailformed")

				hr.ContentLength = 0
			}
		}

		current += 2
		if (*data)[current] == 0x0D {
			current += 2

			if hr.ContentLength != 0 {
				hr.Body = string((*data)[current : current+hr.ContentLength])
			}

			isProcessingFinished = true
		}
	}

	return hr
}

func (hq *HttpRequest) String() string {
	return fmt.Sprintf("Method: %s\nContentLength: %d\nBody:\n%s",
		hq.Method, hq.ContentLength, hq.Body)
}

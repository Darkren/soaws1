package telnetreqparser

import (
	"soawstest/internal/parsers/utility"
)

// TelnetRequestParser is an implementation of reqparser.RequestParser used
// for parsing telnet client requests
type TelnetRequestParser struct {
	payload []byte
}

// New is used to create and initialize TelnetRequestParser instance.
// Inits its payload buffer with nil slice
func New() *TelnetRequestParser {
	return &TelnetRequestParser{[]byte{}}
}

// Next is used to fetch the next ready to be processed request.
// Uses 0x0A (\n) as a requests delimiter
func (trp *TelnetRequestParser) Next() []byte {
	current := utility.Seek(&trp.payload, 0, 0x0A)

	if current == -1 {
		//	in case of win client
		return nil
	}

	req := trp.payload[:current+1]

	trp.payload = trp.payload[len(req):]

	return req
}

// Append is used to append new portion of data to the local payload buffer
func (trp *TelnetRequestParser) Append(data *[]byte) {
	trp.payload = append(trp.payload, (*data)...)
}

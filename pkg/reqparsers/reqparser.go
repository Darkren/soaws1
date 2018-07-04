package reqparsers

// RequestParser is used to parse requests
type RequestParser interface {
	Next() []byte
	Append(data *[]byte)
}

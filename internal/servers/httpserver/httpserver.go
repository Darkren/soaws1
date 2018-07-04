package httpserver

import (
	"fmt"
	"io"
	"net"
	"soawstest/internal/models/requests"
	"soawstest/internal/parsers/httpreqparser"
	"time"
)

const bufferSize = 4096

type HttpServer struct {
	buffer         []byte
	timeoutSeconds int
}

// New is used to create and initialize the server. Inits its local buffer with
// the slice of bufSize size
func New(bufSize int, timeoutSecs int) *HttpServer {
	return &HttpServer{make([]byte, bufSize), timeoutSecs}
}

// StartListening is used to bind to port and start accepting requests
func (s *HttpServer) StartListening(port int) {
	s.buffer = make([]byte, bufferSize)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		fmt.Println("Client connected")

		go s.handleClient(conn)
	}
}

func (s *HttpServer) handleClient(conn net.Conn) {
	conn.SetDeadline(time.Now().Add(time.Duration(s.timeoutSeconds * int(time.Minute))))

	reqParser := httpreqparser.New()

	defer conn.Close()

	for {
		tmp := make([]byte, 128)
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error: ", err)
			}

			fmt.Println("Client disconnected")

			break
		}

		toAppend := tmp[:n]
		reqParser.Append(&toAppend)

		if reqBytes := reqParser.Next(); reqBytes != nil {
			req := &requests.HttpRequest{}
			req = req.New(&reqBytes)
			fmt.Println("Got from client: ", req)

			var resp string
			if req.ContentLength != 0 {
				resp = req.Body
			} else {
				resp = "Hello"

				req.ContentLength = len([]byte(resp))
			}

			conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Length: %d\r\n\r\n%s",
				req.ContentLength, resp)))
		}
	}
}

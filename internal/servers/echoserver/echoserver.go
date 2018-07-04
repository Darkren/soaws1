package echoserver

import (
	"fmt"
	"io"
	"net"
	"soawstest/internal/parsers/telnetreqparser"
)

const bufferSize = 4096

// EchoServer is a tcp server which echoes all requests back to client
type EchoServer struct {
	buffer []byte
}

// New is used to create and initialize the server. Inits its local buffer with
// the slice of bufSize size
func New(bufSize int) *EchoServer {
	return &EchoServer{make([]byte, bufSize)}
}

// StartListening is used to bind to port and start accepting requests
func (s *EchoServer) StartListening(port int) {
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

// handleClient get requests from clients, parses and sends back
func (s *EchoServer) handleClient(conn net.Conn) {
	reqParser := telnetreqparser.New()

	defer conn.Close()

	tmp := make([]byte, 128)
	for {
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

		if cmd := reqParser.Next(); cmd != nil {
			fmt.Println("Got from client: ", string(cmd))

			conn.Write(cmd)
		}
	}
}

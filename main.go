package main

import (
	"soawstest/internal/servers/echoserver"
	"soawstest/pkg/servers"
)

const bufSize = 128

var server servers.Server

func main() {
	server := echoserver.New(bufSize)

	//server := httpserver.New(bufSize, 10)

	server.StartListening(8089)
}

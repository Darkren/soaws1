package servers

type Server interface {
	StartListening(port int)
}

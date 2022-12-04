package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	fmt.Println("server running on port", s.listenAddr)

	go s.acceptLoop()
	s.ln = ln

	<-s.quitch
	close(s.msgch)

	return nil
}

// Accepting new connections
func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()

		if err != nil {
			fmt.Println("accept error:", err)
			continue // Returning the error would stop the loop.
		}
		// New go routine for the read loop.
		go s.readLoop(conn)

	}
}

// Reading the incoming data into the buffer.
func (s *Server) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		num, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		sourceAddr := conn.RemoteAddr().String()

		s.msgch <- Message{
			from:    sourceAddr,
			payload: buf[:num]}

		conn.Write([]byte(fmt.Sprintf("Received, %s!\n", sourceAddr)))
	}
}

func main() {
	listenAddr := ":3000"
	server := NewServer(listenAddr)

	go func() {
		for msg := range server.msgch {
			fmt.Printf("New Message from %s: %s\n", msg.from, string(msg.payload))
		}
	}()

	log.Fatal(server.Start())
}

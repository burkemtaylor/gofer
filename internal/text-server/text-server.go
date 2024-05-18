package textserver

import (
	"fmt"
	"net"
)

type Message struct {
	from    string
	payload []byte
}

func (m Message) GetFrom() string {
	return m.from
}

func (m Message) GetPayload() []byte {
	return m.payload
}

type TextServer struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message
}

func NewTextServer(listenAddr string) *TextServer {
	return &TextServer{
		listenAddr: listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
	}
}

func (s *TextServer) GetMsgCh() chan Message {
	return s.msgch
}

func (s *TextServer) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	close(s.msgch)

	return nil

}

func (s *TextServer) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		fmt.Println("new connection:", conn.RemoteAddr())

		go s.readLoop(conn)
	}
}

func (s *TextServer) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error:", err)
			continue
		}

		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}

		conn.Write([]byte("thank you for your message!"))
	}
}

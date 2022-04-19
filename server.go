package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var (
	conns   []net.Conn
	IP      []string
	connCh  = make(chan net.Conn)
	closeCh = make(chan net.Conn)
	msgCh   = make(chan string)
)

func main() {
	server, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Fatal(err)
			}
			ip_conn := conn.RemoteAddr()
			segments := ip_conn.String()
			for i := range segments {
				if segments[i] == ':' {
					segments = segments[:i]
					break
				}
			}
			IP = append(IP, segments)
			conns = append(conns, conn)
			connCh <- conn
		}
	}()

	for {
		select {
		case conn := <-connCh:
			go onMessage(conn)

		case msg := <-msgCh:
			fmt.Print(msg)

		case conn := <-closeCh:
			fmt.Println("client exit")
			removeConn(conn)
		}
	}
}

func removeConn(conn net.Conn) {
	var i int
	for i = range conns {
		if conns[i] == conn {
			break
		}
	}

	// i = ?  1 2 3 4
	conns = append(conns[:i], conns[i+1:]...)
	IP = append(IP[:i], IP[i+1:]...)
}

func publishMsg(conn net.Conn, msg string, person_IP string) {
	for i := range conns {
		if conns[i] != conn && IP[i] == person_IP {
			conns[i].Write([]byte(msg))
		}
	}
}

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')

		if err != nil {
			break
		}
		/*	for i := range conns {
			if conns[i] == conn {
				msg = msg + IP[i]
			}
		}*/
		// fmt.Println(msg)
		var person_IP string
		for i := range msg {
			if msg[i] == '+' {
				person_IP = msg[i+1 : len(msg)-1]
				msg = msg[:i] + "\n"
				break
			}
		}
		msgCh <- msg
		//	time.Sleep(10 * time.Second)
		publishMsg(conn, msg, person_IP)

	}

	closeCh <- conn
}

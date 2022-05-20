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
	name_IP = make(map[string]string)
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
			IP_conn := conn.RemoteAddr().String()
			userName, _ := name_IP[IP_conn]
			fmt.Printf("client has: IP = %s ; Name = %s exit!\n", IP_conn, userName)
			removeConn(conn) //Remove conn and IP of conn
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
		IP_string := conn.RemoteAddr().String()
		var person_Name string
		for i := range msg {
			if msg[i] == '+' {
				person_Name = msg[:i]
				msg = msg[i+1:]
				break
			}
		}
		for i := range msg {
			if msg[i] == ':' {
				name_IP[IP_string] = msg[:i]
				break
			}
		}
		fmt.Println(name_IP[IP_string])
		msgCh <- msg
		var person_IP string
		for user_IP, user_Name := range name_IP {
			if user_Name == person_Name {
				person_IP = user_IP
			}
		}
		publishMsg(conn, msg, person_IP)

	}

	closeCh <- conn
}

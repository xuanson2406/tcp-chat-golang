package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
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
			fmt.Printf("client has: IP = %s ; Name = %s exit! (%s) \n", IP_conn, userName, time.Now().Format("2006-01-02 15:04:05"))
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
func publishMsg(conn net.Conn, msg string, person_Name string) {
	var person_IP string
	for user_IP, user_Name := range name_IP {
		if user_Name == person_Name {
			person_IP = user_IP
			break
		}
	}
	for i := range conns {
		if conns[i] != conn && IP[i] == person_IP {
			conns[i].Write([]byte(msg))
		}
	}
}
func Reader(conn net.Conn, init_flag *int) string {
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString('\n')
	if err != nil {
		*init_flag = -1
	}
	return msg
}
func onMessage(conn net.Conn) {
	var init_flag int = 0
	for {
		if init_flag == -1 {
			break
		}
		if init_flag == 0 {
			msg := Reader(conn, &init_flag)
			IP_string := conn.RemoteAddr().String()
			name_IP[IP_string] = msg[:len(msg)-1]
			fmt.Printf("client has: IP = %s ; Name = %s --> init socket to the server (%s)\n", IP_string, name_IP[IP_string], time.Now().Format("2006-01-02 15:04:05"))
			init_flag++
		} else {
			msg := Reader(conn, &init_flag)
			var person_Name string
			for i := range msg {
				if msg[i] == '+' {
					person_Name = msg[:i]
					msg = msg[i+1:]
					break
				}
			}
			msgCh <- msg
			publishMsg(conn, msg, person_Name)
		}
	}
	closeCh <- conn
}

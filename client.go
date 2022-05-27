package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func onMessage(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, _ := reader.ReadString('\n')

		fmt.Print(msg)
	}
}

func main() {
	connection, err := net.Dial("tcp", os.Getenv("IP"))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("your name:")
	nameReader := bufio.NewReader(os.Stdin)
	nameInput, _ := nameReader.ReadString('\n')
	nameInput = nameInput[:len(nameInput)-1]
	fmt.Print("Who do you want to chat with ?")
	personReader := bufio.NewReader(os.Stdin)
	PersonInput, _ := personReader.ReadString('\n')
	PersonInput = PersonInput[:len(PersonInput)-1]

	fmt.Println("********** WELCOME TO SOCKET TCP CHAT  **********")
	fmt.Println("*************************************************")
	go onMessage(connection)
	var name string
	name = fmt.Sprintf("%s\n", nameInput)
	connection.Write([]byte(name))
	for {
		msgReader := bufio.NewReader(os.Stdin)
		msg, err := msgReader.ReadString('\n')
		if err != nil {
			break
		}

		msg = fmt.Sprintf("%s+%s: %s\n", PersonInput, nameInput,
			msg[:len(msg)-1])

		connection.Write([]byte(msg))
	}

	connection.Close()
}

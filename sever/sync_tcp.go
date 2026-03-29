package sever

import (
	"fmt"
	"log"
	"net"
	"strconv"

	config "github.com/anujkaushik1/GoDis/Config"
	"github.com/anujkaushik1/GoDis/core"
)

func readCommand(client net.Conn) (*core.RedisCmd, error) {
	buffer := make([]byte, 1024)
	n, err := client.Read(buffer)
	if err != nil {
		return nil, err
	}

	tokens, err := core.Decode(buffer[:n])
	fmt.Println(tokens)

	return nil, nil

}

func respond(client net.Conn, cmd string) error {
	cmd += ""
	_, err := client.Write([]byte(cmd))
	if err != nil {
		return err
	}

	return nil
}

func RunTcpServer() {
	host := config.App.Host
	port := config.App.Port
	log.Println("Starting TCP server on: ", host, ":", port)

	noOfClients := 0

	listener, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal("Failed to start TCP server:", err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		noOfClients++

		log.Println("Total Connected Cliens = ", noOfClients)

		for {
			_, err := readCommand(client)
			if err != nil {
				client.Close()
				noOfClients--
				log.Println("client disconnected")
				break

			}

			client.Write([]byte("+OK\r\n"))
			// err = respond(client, cmd)
			// if err != nil {
			// 	client.Close()
			// 	noOfClients--
			// 	log.Println("client disconnected2222")
			// 	break

			// }

		}

	}

}

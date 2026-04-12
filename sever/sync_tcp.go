package sever

import (
	"github.com/anujkaushik1/GoDis/core"
)

func Respond(client core.FileDescriptor, redisCmds *core.RedisCmds) error {
	var buf []byte
	for _, redisCmd := range *redisCmds {
		buf = append(buf, core.Eval(redisCmd)...)
	}
	_, err := client.Write(buf)
	return err
}

// func RunTcpServer() {
// 	host := config.App.Host
// 	port := config.App.Port
// 	log.Println("Starting TCP server on: ", host, ":", port)

// 	noOfClients := 0

// 	listener, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
// 	if err != nil {
// 		log.Fatal("Failed to start TCP server:", err)
// 	}

// 	// for {
// 	// 	client, err := listener.Accept()
// 	// 	if err != nil {
// 	// 		panic(err)
// 	// 	}

// 	// 	noOfClients++

// 	// 	log.Println("Total Connected Cliens = ", noOfClients)

// 	// 	for {
// 	// 		redisCmd, err := ReadCommand(client)
// 	// 		if err != nil {
// 	// 			client.Close()
// 	// 			noOfClients--
// 	// 			log.Println("client disconnected")
// 	// 			break

// 	// 		}

// 	// 		err = respond(client, redisCmd)
// 	// 		if err != nil {
// 	// 			client.Close()
// 	// 			noOfClients--
// 	// 			log.Println("client disconnected2222")
// 	// 			break

// 	// 		}

// 	// 	}

// 	// }

// }

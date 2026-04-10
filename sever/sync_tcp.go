package sever

import (
	"fmt"

	"github.com/anujkaushik1/GoDis/core"
)

func ReadCommand(client core.FileDescriptor) (*core.RedisCmds, error) {
	buffer := make([]byte, 1024)
	n, err := client.Read(buffer)
	if err != nil {
		return nil, err
	}

	tokens := core.DecodeAndFlatten(buffer[:n])
	if len(tokens) == 0 {
		return nil, fmt.Errorf("unable to decode command")
	}

	redisCmds := make(core.RedisCmds, 0)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		redisCmd := &core.RedisCmd{
			Cmd:  token[0],
			Args: token[1:],
		}

		redisCmds = append(redisCmds, redisCmd)
	}

	return &redisCmds, nil

}

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

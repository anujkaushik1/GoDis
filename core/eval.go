package core

import (
	"strconv"
)

func EvalAndRespond(client FileDescriptor, redisCmd *RedisCmd) error {
	cmd := redisCmd.Cmd
	args := redisCmd.Args

	if cmd == "COMMAND" {
		_, err := client.Write([]byte("+OK\r\n"))
		return err
	}
	if cmd == "PING" {
		if len(args) == 0 {
			_, err := client.Write([]byte("+PONG\r\n"))
			return err
		}
		if len(args) == 1 {
			_, err := client.Write([]byte("$" + strconv.Itoa(len(args[0])) + "\r\n" + args[0] + "\r\n"))
			return err
		}
		if len(args) > 1 {
			_, err := client.Write([]byte("-ERR wrong number of arguments for 'ping' command\r\n"))
			return err
		}
		return nil

	}

	return nil

}

package core

import (
	"strconv"
	"time"
)

func EvalPING(client FileDescriptor, args []string) error {
	if len(args) == 0 {
		_, err := client.Write([]byte("+PONG\r\n"))
		return err
	}
	if len(args) == 1 {
		_, err := client.Write([]byte("$" + strconv.Itoa(len(args[0])) + "\r\n" + args[0] + "\r\n"))
		return err
	}
	if len(args) > 1 {
		return EvalErr(client, "ERR wrong number of arguments for 'ping' command")
	}
	return nil
}

func EvalErr(client FileDescriptor, msg string) error {
	_, err := client.Write([]byte("-" + msg + "\r\n"))
	return err
}

func EvalString(client FileDescriptor, msg string) error {
	_, err := client.Write([]byte("+" + msg + "\r\n"))
	return err
}

func EvalNil(client FileDescriptor) error {
	_, err := client.Write([]byte("$-1\r\n"))
	return err
}

func EvalInteger(client FileDescriptor, val int64) error {
	_, err := client.Write([]byte(":" + strconv.FormatInt(val, 10) + "\r\n"))
	return err
}

func EvalSET(client FileDescriptor, args []string) error {

	if len(args) <= 1 {
		return EvalErr(client, "ERR wrong number of arguments for 'set' command")
	}

	key, value := args[0], args[1]

	expiresIn := -1
	for i := 2; i < len(args); i++ {
		if args[i] == "EX" || args[i] == "ex" {
			newIdx := i + 1
			if newIdx == len(args) {
				return EvalErr(client, "ERR syntax error")
			}

			intVal, err := strconv.Atoi(args[newIdx])
			if err != nil {
				return EvalErr(client, "ERR value is not an integer or out of range")
			}
			expiresIn = intVal
		}
	}

	Set(key, value, int64(expiresIn))
	return EvalString(client, "OK")

}

func EvalGET(client FileDescriptor, args []string) error {

	if len(args) > 1 {
		return EvalErr(client, "ERR wrong number of arguments for 'get' command")
	}

	storeObj := Get(args[0])

	if storeObj == nil {
		return EvalNil(client)
	}

	expiresIn := storeObj.ExpiresAt
	value := storeObj.Value

	if expiresIn > 0 && time.Now().UnixMilli() > expiresIn {
		return EvalNil(client)
	}

	strValue, ok := value.(string)
	if !ok {
		return EvalNil(client)
	}
	return EvalString(client, strValue)

}

func EvalTTL(client FileDescriptor, args []string) error {
	if len(args) > 1 {
		return EvalErr(client, "ERR wrong number of arguments for 'ttl' command")
	}

	storeObj := Get(args[0])

	if storeObj == nil {
		return EvalInteger(client, -2)
	}

	expiresIn := storeObj.ExpiresAt

	if expiresIn < 0 {
		return EvalInteger(client, -1)
	}

	if time.Now().UnixMilli() > expiresIn {
		return EvalInteger(client, -2)
	}

	ttl := (expiresIn - time.Now().UnixMilli()) / 1000
	return EvalInteger(client, ttl)

}

func EvalExpire(client FileDescriptor, args []string) error {

	if len(args) <= 1 {
		return EvalErr(client, "ERR wrong number of arguments for 'expire' command")
	}

	storeObj := Get(args[0])

	if storeObj == nil {
		return EvalInteger(client, 0)
	}
	value := storeObj.Value
	expiresIn, err := strconv.Atoi(args[1])
	if err != nil {
		return EvalErr(client, "ERR value is not an integer or out of range")
	}

	Set(args[0], value, int64(expiresIn))
	return EvalInteger(client, 1)

}
func EvalAndRespond(client FileDescriptor, redisCmd *RedisCmd) error {
	cmd := redisCmd.Cmd
	args := redisCmd.Args

	if cmd == "COMMAND" {
		return EvalString(client, "OK")
	}
	if cmd == "PING" {
		return EvalPING(client, args)
	}

	if cmd == "SET" {
		return EvalSET(client, args)
	}

	if cmd == "GET" {
		return EvalGET(client, args)
	}

	if cmd == "TTL" {
		return EvalTTL(client, args)
	}

	if cmd == "EXPIRE" {
		return EvalExpire(client, args)
	}

	return nil

}

package core

import (
	"strconv"
	"time"
)

func EvalPING(args []string) []byte {
	if len(args) == 0 {
		return []byte("+PONG\r\n")
	}
	if len(args) == 1 {
		return []byte("$" + strconv.Itoa(len(args[0])) + "\r\n" + args[0] + "\r\n")
	}
	if len(args) > 1 {
		return EvalErr("ERR wrong number of arguments for 'ping' command")
	}
	return nil
}

func EvalErr(msg string) []byte {
	return []byte("-" + msg + "\r\n")
}

func EvalString(msg string) []byte {
	return []byte("+" + msg + "\r\n")
}

func EvalNil() []byte {
	return []byte("$-1\r\n")
}

func EvalInteger(val int64) []byte {
	return []byte(":" + strconv.FormatInt(val, 10) + "\r\n")
}

func EvalSET(args []string) []byte {

	if len(args) <= 1 {
		return EvalErr("ERR wrong number of arguments for 'set' command")
	}

	key, value := args[0], args[1]

	expiresIn := -1
	for i := 2; i < len(args); i++ {
		if args[i] == "EX" || args[i] == "ex" {
			newIdx := i + 1
			if newIdx == len(args) {
				return EvalErr("ERR syntax error")
			}

			intVal, err := strconv.Atoi(args[newIdx])
			if err != nil {
				return EvalErr("ERR value is not an integer or out of range")
			}
			expiresIn = intVal
		}
	}

	Set(key, value, int64(expiresIn))
	return EvalString("OK")

}

func EvalGET(args []string) []byte {

	if len(args) > 1 {
		return EvalErr("ERR wrong number of arguments for 'get' command")
	}

	storeObj := Get(args[0])

	if storeObj == nil {
		return EvalNil()
	}

	expiresIn := storeObj.ExpiresAt
	value := storeObj.Value

	if expiresIn > 0 && time.Now().UnixMilli() > expiresIn {
		Del(args[0])
		return EvalNil()
	}

	strValue, ok := value.(string)
	if !ok {
		return EvalNil()
	}
	return EvalString(strValue)

}

func EvalTTL(args []string) []byte {
	if len(args) > 1 {
		return EvalErr("ERR wrong number of arguments for 'ttl' command")
	}

	storeObj := Get(args[0])

	if storeObj == nil {
		return EvalInteger(-2)
	}

	expiresIn := storeObj.ExpiresAt

	if expiresIn < 0 {
		return EvalInteger(-1)
	}

	if time.Now().UnixMilli() > expiresIn {
		return EvalInteger(-2)
	}

	ttl := (expiresIn - time.Now().UnixMilli()) / 1000
	return EvalInteger(ttl)

}

func EvalExpire(args []string) []byte {

	if len(args) <= 1 {
		return EvalErr("ERR wrong number of arguments for 'expire' command")
	}

	storeObj := Get(args[0])

	if storeObj == nil {
		return EvalInteger(0)
	}
	value := storeObj.Value
	expiresIn, err := strconv.Atoi(args[1])
	if err != nil {
		return EvalErr("ERR value is not an integer or out of range")
	}

	Set(args[0], value, int64(expiresIn))
	return EvalInteger(1)

}

func EvalDel(args []string) []byte {

	if len(args) == 0 {
		return EvalErr("ERR wrong number of arguments for 'del' command")
	}

	deletedCount := 0

	for key := 0; key < len(args); key++ {
		if Del(args[key]) {
			deletedCount++
		}
	}

	return EvalInteger(int64(deletedCount))
}

func Eval(redisCmd *RedisCmd) []byte {
	cmd := redisCmd.Cmd
	args := redisCmd.Args

	if cmd == "COMMAND" {
		return EvalString("OK")
	}
	if cmd == "PING" {
		return EvalPING(args)
	}

	if cmd == "SET" {
		return EvalSET(args)
	}

	if cmd == "GET" {
		return EvalGET(args)
	}

	if cmd == "TTL" {
		return EvalTTL(args)
	}

	if cmd == "EXPIRE" {
		return EvalExpire(args)
	}

	if cmd == "DEL" {
		return EvalDel(args)
	}

	return nil

}

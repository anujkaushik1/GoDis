package core

import "fmt"

func ReadCommand(client FileDescriptor) (*RedisCmds, error) {
	buffer := make([]byte, 1024)
	n, err := client.Read(buffer)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}

	tokens := DecodeAndFlatten(buffer[:n])
	if len(tokens) == 0 {
		return nil, fmt.Errorf("unable to decode command")
	}

	redisCmds := make(RedisCmds, 0)

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		redisCmd := &RedisCmd{
			Cmd:  token[0],
			Args: token[1:],
		}

		redisCmds = append(redisCmds, redisCmd)
	}

	return &redisCmds, nil
}

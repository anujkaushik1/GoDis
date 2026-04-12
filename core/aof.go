package core

import (
	"fmt"
	"os"
)

func DumpAllAof() {
	file, err := os.OpenFile("appendonly.aof", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	for key, val := range store {
		value := fmt.Sprintf("%v", val.Value)
		cmd := fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n", len(key), key, len(value), value)
		file.WriteString(cmd)
	}

	fmt.Println("Aof file rewrite complete")

}

func AppendToAof(cmd string) {
	file, err := os.OpenFile("appendonly.aof", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	file.WriteString(cmd)
}

func LoadAofFile() {
	fmt.Println("Loading aof file....")
	aofFile, err := os.OpenFile(
		"appendonly.aof",
		os.O_CREATE|os.O_RDONLY,
		0644,
	)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer aofFile.Close()

	fd := FileDescriptor{FD: int(aofFile.Fd())}
	redisCmds, err := ReadCommand(fd)
	if err != nil || redisCmds == nil {
		return
	}

	BulkWriteToStore(redisCmds)
	fmt.Println("Loaded aof file sucessfully")

}

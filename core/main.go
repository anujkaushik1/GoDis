package main

import (
	"errors"
	"fmt"
	"strconv"
)

func readSimpleString(bytes []byte) (string, int, error) {

	pos := 1

	for bytes[pos] != '\r' {
		pos++
	}

	return string(bytes[1:pos]), pos + 2, nil

}

func readError(bytes []byte) (string, int, error) {
	return readSimpleString(bytes)
}

func readInt64(bytes []byte) (int64, int, error) {

	pos := 1

	for bytes[pos] != '\r' {
		pos++
	}

	val, err := strconv.ParseInt(string(bytes[1:pos]), 10, 64)

	if err != nil {
		return 0, 0, err
	}

	return val, pos + 2, nil

}

func readBulkString(bytes []byte) (string, int, error) {
	pos := 1

	for bytes[pos] != '\r' {
		pos++
	}

	pos += 2

	startOfString := pos

	for bytes[pos] != '\r' {
		pos++
	}

	return string(bytes[startOfString:pos]), pos + 2, nil

}

func decodeOne(bytes []byte) (interface{}, int, error) {

	switch bytes[0] {
	case '+':
		return readSimpleString(bytes)

	case '-':
		return readSimpleString(bytes)

	case ':':
		return readInt64(bytes)

	case '$':
		return readBulkString(bytes)

	}

	return nil, 0, nil

}

func decode(bytes []byte) (interface{}, error) {

	if len(bytes) == 0 {
		return nil, errors.New("No data found")
	}

	value, _, err := decodeOne(bytes)

	return value, err

}

func main() {
	val, err := decode([]byte("$5\r\nanujk\r\n"))
	if err != nil {
		fmt.Println("error in decode = ", err.Error())
	}
	fmt.Println("ansss = ", val)
}

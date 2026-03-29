package main

import (
	"errors"
	"fmt"
)

func readSimpleString(bytes []byte) (string, int, error) {

	pos := 1

	for bytes[pos] != '\r' {
		pos++
	}

	return string(bytes[1:pos]), pos, nil

}

func decodeOne(bytes []byte) (interface{}, int, error) {

	switch bytes[0] {
	case '+':
		return readSimpleString(bytes)

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
	val, _ := decode([]byte("+Hhahahelloworld\r\n"))
	fmt.Println(val)
}

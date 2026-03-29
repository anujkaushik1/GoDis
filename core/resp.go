package core

import (
	"errors"
	"strconv"
)

func ReadSimpleString(bytes []byte) (string, int, error) {

	pos := 1

	for bytes[pos] != '\r' {
		pos++
	}

	return string(bytes[1:pos]), pos + 2, nil

}

func ReadError(bytes []byte) (string, int, error) {
	return ReadSimpleString(bytes)
}

func ReadInt64(bytes []byte) (int64, int, error) {

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

func ReadBulkString(bytes []byte) (string, int, error) {
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

func ReadArray(bytes []byte) (interface{}, int, error) {

	pos := 1

	numberOfElements := 0

	for bytes[pos] != '\r' {
		numberOfElements = numberOfElements*10 + int(bytes[pos]-'0')
		pos++
	}

	pos += 2

	arrayElements := make([]interface{}, numberOfElements)
	arrayIdx := 0
	for arrayIdx < len(arrayElements) {
		value, delta, err := DecodeOne(bytes[pos:])
		if err != nil {
			return nil, 0, err
		}
		pos += delta

		arrayElements[arrayIdx] = value
		arrayIdx++

	}

	return arrayElements, pos, nil
}

func DecodeOne(bytes []byte) (interface{}, int, error) {

	switch bytes[0] {
	case '+':
		return ReadSimpleString(bytes)

	case '-':
		return ReadSimpleString(bytes)

	case ':':
		return ReadInt64(bytes)

	case '$':
		return ReadBulkString(bytes)

	case '*':
		return ReadArray(bytes)

	}

	return nil, 0, nil

}

func Decode(bytes []byte) (interface{}, error) {

	if len(bytes) == 0 {
		return nil, errors.New("No data found")
	}

	value, _, err := DecodeOne(bytes)

	return value, err

}

package core

import (
	"errors"
	"strconv"
)

func ReadSimpleString(bytes []byte) (string, int, error) {

	pos := 1

	for pos < len(bytes) && bytes[pos] != '\r' {
		pos++
	}

	if pos >= len(bytes) {
		return "", 0, errors.New("invalid string: missing \\r\\n terminator")
	}

	return string(bytes[1:pos]), pos + 2, nil

}

func ReadError(bytes []byte) (string, int, error) {
	return ReadSimpleString(bytes)
}

func ReadInt64(bytes []byte) (int64, int, error) {

	pos := 1

	for pos < len(bytes) && bytes[pos] != '\r' {
		pos++
	}

	if pos >= len(bytes) {
		return 0, 0, errors.New("invalid integer: missing \\r\\n terminator")
	}

	val, err := strconv.ParseInt(string(bytes[1:pos]), 10, 64)

	if err != nil {
		return 0, 0, err
	}

	return val, pos + 2, nil

}

func ReadBulkString(bytes []byte) (string, int, error) {
	pos := 1

	for pos < len(bytes) && bytes[pos] != '\r' {
		pos++
	}

	if pos >= len(bytes) {
		return "", 0, errors.New("invalid bulk string: missing \\r\\n terminator")
	}

	pos += 2

	startOfString := pos

	for pos < len(bytes) && bytes[pos] != '\r' {
		pos++
	}

	if pos >= len(bytes) {
		return "", 0, errors.New("invalid bulk string: missing \\r\\n terminator")
	}

	return string(bytes[startOfString:pos]), pos + 2, nil

}

func ReadArray(bytes []byte) (any, int, error) {

	pos := 1

	numberOfElements := 0

	for pos < len(bytes) && bytes[pos] != '\r' {
		numberOfElements = numberOfElements*10 + int(bytes[pos]-'0')
		pos++
	}

	if pos >= len(bytes) {
		return nil, 0, errors.New("invalid array: missing \\r\\n terminator")
	}

	pos += 2

	arrayElements := make([]any, numberOfElements)
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

func DecodeOne(bytes []byte) (any, int, error) {

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

func Decode(bytes []byte) (any, error) {

	if len(bytes) == 0 {
		return nil, errors.New("No data found")
	}

	value, _, err := DecodeOne(bytes)

	return value, err

}

func DecodeAndFlatten(bytes []byte) []string {
	v, err := Decode(bytes)
	if err != nil {
		return nil
	}
	return flatten(v)
}

func flatten(v any) []string {
	result := make([]string, 0)

	switch val := v.(type) {
	case string:
		result = append(result, val)
	case int64:
		result = append(result, strconv.FormatInt(val, 10))
	case []any:
		for _, item := range val {
			result = append(result, flatten(item)...)
		}
	}

	return result
}

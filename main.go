package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		resp := &Resp{reader: bufio.NewReader(conn)}
		value, err := resp.ReadFromBuffer()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(value)

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}

type Value struct {
	typ   string //typ is used to determine the data type carried by the value.
	str   string //str holds the value of the string received from the simple strings.
	num   int    //num holds the value of the integer received from the integers.
	bulk  string //bulk is used to store the string received from the bulk strings.
	array []Value
}

type Writer struct {
	writer io.Writer
}

type Resp struct {
	reader *bufio.Reader
}

// readLine reads and removes the bytes from the buffer.
// This functiona also ensures that we remove the CRLF from the buffer.
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' { // We endup at this if when the reader stops at the \n
			break
		}
	}
	return line[:len(line)-2], n, nil // So we must take all the bytes before \n
}

func (r *Resp) readInteger() (number int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	unifedLine := string(line) // []byte{ '1', '2', '3' }, becomes"123"

	i64, err := strconv.ParseInt(unifedLine, 10, 64) // 10 indicates that is a decimal base. 64 is the size of the integer
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

func (r *Resp) ReadFromBuffer() (value Value, err error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		fmt.Println("Error reading byte:", err)
		return
	}

	fmt.Println("First byte of the sequence: " + string(_type))

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()

	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	value := Value{}
	value.typ = "array"

	arraySize, _, err := r.readInteger()
	if err != nil {
		return value, err
	}

	value.array = make([]Value, arraySize)

	for i := 0; i < arraySize; i++ {
		val, err := r.ReadFromBuffer()
		if err != nil {
			return value, err
		}

		// add parsed value to array
		value.array[i] = val
	}

	return value, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	lengthOfBulk, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, lengthOfBulk)

	r.reader.Read(bulk) // This reads the buffer exactly until lengthOfBulk and saves the bytes read on the slice of bytes

	v.bulk = string(bulk) // Now with all bytes on the slice, we can create the string. Example: []byte{'H', 'e', 'l', 'l', 'o'} => "Hello"

	r.readLine() // And we finally read the line to remove the bytes and CRLF from the buffer.

	return v, nil
}

func (v Value) marshallString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r')
	bytes = append(bytes, '\n')

	return bytes
}

func (v Value) marshallBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, byte(v.num))

	return bytes
}

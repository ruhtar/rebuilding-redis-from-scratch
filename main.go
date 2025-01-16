package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

func main() {
	port := "6379"
	fmt.Println("Listening on port: " + port)

	// Create a new server
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return
	}

	playingWithByteReads()

	// Listen for connections
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {
		buf := make([]byte, 1024)

		// read message from client
		_, err = conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from client: ", err.Error())
			os.Exit(1)
		}

		conn.Write([]byte("+OK\r\n"))
	}
}

func playingWithByteReads() {
	input := "$5\r\nAhmed\r\n"
	reader := bufio.Reader(*bufio.NewReader(strings.NewReader(input)))

	firstByte, _ := reader.ReadByte()
	if firstByte != '$' {
		fmt.Println("Invalid type, expecting bulk strings only")
		os.Exit(1)
	}

	secondByte, _ := reader.ReadByte()
	size, _ := strconv.ParseInt(string(secondByte), 10, 64)

	fmt.Println(size)

	// consume /r/n
	reader.ReadByte()
	reader.ReadByte()

	name := make([]byte, size)
	reader.Read(name)

	fmt.Println(string(name))
}

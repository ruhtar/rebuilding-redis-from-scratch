package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr := "127.0.0.1:7777"

	// Conecta ao servidor
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server. Type a message and press Enter.")

	// Loop para enviar mensagens ao servidor
	for {
		fmt.Print("Enter message: ")
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		// Envia a mensagem ao servidor
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		// Lê a resposta do servidor
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from server:", err)
			return
		}

		fmt.Println("Server response:", string(buf[:n]))
	}
}

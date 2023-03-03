package main

import (
	"fmt"
	"net"
)

func main() {
	host := ""      // Host IP address, empty string means the server will bind to all available interfaces
	port := "12345" // Arbitrary port number, can be any number above 1024
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on port", port)

	for {
		conn, err := listener.Accept() // Wait for a client to connect
		if err != nil {
			fmt.Println("Error accepting:", err.Error())
			continue
		}
		fmt.Println("Client", conn.RemoteAddr(), "connected")

		go handleClient(conn) // Handle client connection in a separate goroutine
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf) // Receive data from client
		if err != nil {          // If connection has been closed
			fmt.Println("Client", conn.RemoteAddr(), "disconnected")
			return
		}
		data := buf[:n]
		fmt.Println("Received from client", conn.RemoteAddr(), ":", string(data))

		go func() {
			_, err := conn.Write(data) // Send data back to client
			if err != nil {
				fmt.Println("Error sending data to client:", err.Error())
				return
			}
		}()
	}
}

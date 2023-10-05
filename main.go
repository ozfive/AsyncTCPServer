package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	host := ""      // Host IP address, empty string means the server will bind to all available interfaces
	port := "12345" // Arbitrary port number, can be any number above 1024
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		log.Fatalf("Error listening: %v", err)
	}
	defer listener.Close()
	log.Printf("Server listening on port %s", port)

	// Handle termination signals gracefully
	setupSignalHandler()

	for {
		conn, err := listener.Accept() // Wait for a client to connect
		if err != nil {
			log.Printf("Error accepting: %v", err)
			continue
		}
		log.Printf("Client %s connected", conn.RemoteAddr())

		go handleClient(conn) // Handle client connection in a separate goroutine
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf) // Receive data from client
		if err != nil {          // If connection has been closed
			log.Printf("Client %s disconnected", conn.RemoteAddr())
			return
		}
		data := buf[:n]
		log.Printf("Received from client %s: %s", conn.RemoteAddr(), string(data))

		go func() {
			_, err := conn.Write(data) // Send data back to client
			if err != nil {
				log.Printf("Error sending data to client: %v", err)
				return
			}
		}()
	}
}

func setupSignalHandler() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("\nReceived termination signal. Cleaning up resources...")
		os.Exit(0)
	}()
}

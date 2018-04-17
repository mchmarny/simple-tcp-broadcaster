package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"unicode"
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func swapCase(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case unicode.IsLower(r):
			return unicode.ToUpper(r)
		case unicode.IsUpper(r):
			return unicode.ToLower(r)
		}
		return r
	}, s)
}

func handleServerConnection(c net.Conn, i int) {
	for {
		// scan message
		scanner := bufio.NewScanner(c)
		for scanner.Scan() {
			msg := scanner.Text()
			log.Printf("Client %v sends: %v", i, msg)
			msgNew := swapCase(reverse((msg)))
			c.Write([]byte(msgNew + "\n"))
			log.Printf("Client %v receives: %v", i, msgNew)
			fmt.Println("---")
		}
		if errRead := scanner.Err(); errRead != nil {
			log.Printf("Client %v disconnected...", i)
			fmt.Println("---")
			return
		}
	}
}

// StartServer start TCP server on a given port
func StartServer(port int) error {
	log.Println("Launching sever...")
	// listen on all interfaces
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	log.Printf("Server launched on port:%d, waiting for clients...\n", port)

	i := 0
	for {
		// accept connection on port
		c, err := ln.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
		}
		i++
		log.Printf("Client %v connected...", i)
		fmt.Println("---")
		// handle the connection
		go handleServerConnection(c, i)
	}

}

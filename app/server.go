package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!!!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {

	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		msg := string(buf[:n])
		fmt.Printf("%s message received", msg)

		response := "+PONG\r\n"
		_, err = conn.Write([]byte(response))
		if err != nil {
			fmt.Println("Error writing:", err)
			return
		}
	}

}

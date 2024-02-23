package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:"+serverContext.cfg.Port)
	if err != nil {
		fmt.Println("Failed to bind to port: ", serverContext.cfg.Port)
		os.Exit(1)
	}
	defer listener.Close()

	var replicationChannel chan []byte

	if serverContext.replicationStore.IsReplica() {
		masterConn, err := net.Dial("tcp", serverContext.replicationStore.MasterAddress)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		handleHandshake(masterConn)
		go handleConnection(masterConn, true, nil)
	} else {
		replicationChannel = make(chan []byte)
		go handleSyncWithReplicas(replicationChannel)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn, false, replicationChannel)
	}
}

func handleConnection(conn net.Conn, isPersistentConn bool, replicationChannel chan []byte) {
	writer := bufio.NewWriter(conn)
	buf := make([]byte, 1024)

	defer func() {
		if !isPersistentConn && !serverContext.replicationStore.IsReplicaClient(conn) {
			conn.Close()
		}
	}()

	for {
		bytesRead, err := conn.Read(buf)
		if err != nil {
			break
		}

		bytesData := buf[:bytesRead]

		processedResult := serverContext.cmdProcessor.ProcessCmd(bytesData, conn)

		for _, item := range processedResult {
			if len(item.Answer) > 0 {
				writer.Write([]byte(item.Answer))
			}

			if replicationChannel != nil && item.IsDuplicate {
				replicationChannel <- item.BytesInput
			}
		}

		writer.Flush()

	}

}

func handleSyncWithReplicas(replicationChannel chan []byte) {
	for {
		serverContext.replicationStore.PopulateCmdToReplicas(<-replicationChannel)
	}
}

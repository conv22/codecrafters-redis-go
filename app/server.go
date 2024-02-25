package main

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/cmds"
	"github.com/codecrafters-io/redis-starter-go/app/handshake"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
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
		handshake.NewMasterHandshake(serverContext.cfg, serverContext.inMemoryStorage).HandleHandshake(masterConn)
		go handleConnection(masterConn, nil, true)
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
		go handleConnection(conn, replicationChannel, false)
	}
}

func handleConnection(conn net.Conn, replicationChannel chan []byte, keepAlive bool) {
	writer := bufio.NewWriter(conn)
	buf := make([]byte, 1024)

	cmdProcessor := cmds.NewRespCmdProcessor(serverContext.inMemoryStorage, serverContext.cfg, serverContext.replicationStore, conn)

	defer func() {
		if !keepAlive && !serverContext.replicationStore.IsReplicaClient(conn) {
			conn.Close()
		}
	}()

	for {
		bytesRead, err := conn.Read(buf)
		if err != nil {
			break
		}

		bytesData := buf[:bytesRead]

		parsed, err := resp.HandleParse(string(bytesData))

		if err != nil {
			continue
		}

		if serverContext.replicationStore.IsMaster() && len(parsed) > 0 && len(parsed[0]) > 0 && parsed[0][0].Value == handshake.HANDSHAKE_CMD_REPLCONF {
			err := handshake.NewClientHandshake(serverContext.cfg, serverContext.inMemoryStorage, serverContext.replicationStore).HandleHandshake(conn, parsed[0])

			if err != nil {
				conn.Close()
				writer.Flush()
				break
			}

		}
		processedResult := cmdProcessor.ProcessCmd(parsed, conn)

		for _, item := range processedResult {
			if len(item.Answer) > 0 {
				writer.Write([]byte(item.Answer))
			}

			if replicationChannel != nil && item.IsPropagate {
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

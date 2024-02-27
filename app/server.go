package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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
		go handleMasterConnection()
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

func handleMasterConnection() {
	masterConn, err := net.Dial("tcp", serverContext.replicationStore.MasterAddress)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	masterHandshake := handshake.NewMasterHandshake(serverContext.cfg, serverContext.inMemoryStorage)
	masterHandshake.HandleHandshake(masterConn)
	handleConnection(masterConn, nil, true)
}

func handleConnection(conn net.Conn, replicationChannel chan []byte, keepAlive bool) {
	writer := bufio.NewWriter(conn)
	buf := make([]byte, 1024)
	cmdProcessor := cmds.NewRespCmdProcessor(serverContext.inMemoryStorage, serverContext.cfg, serverContext.replicationStore, conn)
	processingChannel := make(chan []cmds.ProcessCmdResult)
	go processConnectionCommands(processingChannel, replicationChannel, writer)

	defer conn.Close()
	defer close(replicationChannel)

	for {
		bytesRead, err := conn.Read(buf)
		if err != nil {
			continue
		}

		parsed, err := resp.HandleParse(string(buf[:bytesRead]))

		if err != nil {
			continue
		}

		if isHandshakeStartedRequest(parsed) {
			nextCmds, err := performClientHandshake(conn, parsed)
			if err != nil {
				continue
			}
			if nextCmds != nil {
				parsed = nextCmds
			}
		}

		processingChannel <- cmdProcessor.ProcessCmd(parsed, conn)
	}

}

func performClientHandshake(conn net.Conn, parsed []resp.ParsedCmd) (nextCmds []resp.ParsedCmd, err error) {
	clientHandshake := handshake.NewClientHandshake(serverContext.cfg, serverContext.inMemoryStorage, serverContext.replicationStore)

	return clientHandshake.HandleHandshake(conn, parsed)

}

func processConnectionCommands(output chan []cmds.ProcessCmdResult, replicationChannel chan []byte, writer *bufio.Writer) {
	for outputResult := range output {
		for _, item := range outputResult {
			if replicationChannel != nil && item.IsPropagate {
				replicationChannel <- item.BytesInput
			}

			if len(item.Answer) > 0 {
				writer.Write([]byte(item.Answer))
				writer.Flush()
			}
		}
	}

}
func isHandshakeStartedRequest(parsed []resp.ParsedCmd) bool {
	return serverContext.replicationStore.IsMaster() &&
		len(parsed) > 0 &&
		strings.EqualFold(parsed[0].Value, handshake.HANDSHAKE_CMD_REPLCONF)
}

func handleSyncWithReplicas(replicationChannel chan []byte) {
	for {
		serverContext.replicationStore.PopulateCmdToReplicas(<-replicationChannel)
	}
}

package main

import (
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
			return
		}
		handshake.New(serverContext.cfg, serverContext.inMemoryStorage).HandleHandshake(masterConn)
		go handleConnection(masterConn, nil)
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
		go handleConnection(conn, replicationChannel)
	}
}

func handleConnection(conn net.Conn, replicationChannel chan []byte) {
	cmdProcessor := cmds.NewRespCmdProcessor(serverContext.inMemoryStorage, serverContext.cfg, serverContext.replicationStore, conn)
	processingChannel := make(chan []byte)
	reader := resp.NewReader(conn)

	defer func() {
		conn.Close()
		if replicationChannel != nil {
			close(replicationChannel)
		}
	}()

	go func() {
		for outputResult := range processingChannel {
			conn.Write(outputResult)
		}
	}()

	for {

		parsed, err := reader.HandleRead()

		if err != nil || len(parsed) == 0 {
			continue
		}

		output := cmdProcessor.ProcessCmd(parsed, conn)

		for _, item := range output {
			if replicationChannel != nil && item.IsPropagate {
				replicationChannel <- item.BytesInput
			}

			if len(item.Answer) > 0 {
				processingChannel <- []byte(item.Answer)
			}
		}
	}

}

func handleSyncWithReplicas(replicationChannel chan []byte) {
	for {
		serverContext.replicationStore.PopulateCmdToReplicas(<-replicationChannel)
	}
}

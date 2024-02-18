package cmds

import (
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/replication"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
	"github.com/codecrafters-io/redis-starter-go/app/storage"
)

var RespEncodingConstants = resp.RESP_ENCODING_CONSTANTS

type RespCmdProcessor struct {
	parser      *resp.RespParser
	storage     *storage.StorageCollection
	config      *config.Config
	replication *replication.ReplicationInfo
	connection  net.Conn
}

func NewRespCmdProcessor(p *resp.RespParser, storage *storage.StorageCollection, config *config.Config, replication *replication.ReplicationInfo, conn net.Conn) *RespCmdProcessor {
	return &RespCmdProcessor{
		parser:      p,
		storage:     storage,
		config:      config,
		replication: replication,
		connection:  conn,
	}
}

const (
	CMD_PING        string = "PING"
	CMD_PONG        string = "PONG"
	CMD_ECHO        string = "ECHO"
	CMD_GET         string = "GET"
	CMD_SET         string = "SET"
	CMD_CONFIG      string = "CONFIG"
	CMD_KEYS        string = "KEYS"
	CMD_INFO        string = "INFO"
	CMD_REPLCONF    string = "REPLCONF"
	CMD_PSYNC       string = "PSYNC"
	CMD_OK          string = "OK"
	CMD_FULL_RESYNC string = "FULLRESYNC"
)

type ProcessCmdResult struct {
	Answer      string
	IsDuplicate bool
}

func (processor *RespCmdProcessor) ProcessCmd(line string) []ProcessCmdResult {
	result := []ProcessCmdResult{}
	parsedLines, err := processor.parser.HandleParse(line)

	if err != nil {
		return []ProcessCmdResult{{Answer: processor.parser.HandleEncode(RespEncodingConstants.ERROR, "error parsing the line")}}
	}

	if len(parsedLines) == 0 {
		return []ProcessCmdResult{{Answer: processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not enough arguments")}}
	}

	for _, parsedLine := range parsedLines {

		firstCmd := strings.ToUpper(parsedLine[0].Value)
		cmds := parsedLine[1:]

		switch firstCmd {
		case CMD_PING:
			result = append(result, ProcessCmdResult{Answer: processor.handlePing()})
		case CMD_ECHO:
			result = append(result, ProcessCmdResult{Answer: processor.handleEcho(cmds)})
		case CMD_SET:
			result = append(result, ProcessCmdResult{Answer: processor.handleSet(cmds), IsDuplicate: true})
		case CMD_GET:
			result = append(result, ProcessCmdResult{Answer: processor.handleGet(cmds)})
		case CMD_CONFIG:
			result = append(result, ProcessCmdResult{Answer: processor.handleConfig(cmds)})
		case CMD_KEYS:
			result = append(result, ProcessCmdResult{Answer: processor.handleKeys(cmds)})
		case CMD_INFO:
			result = append(result, ProcessCmdResult{Answer: processor.handleInfo(cmds)})

		case CMD_REPLCONF:
			result = append(result, ProcessCmdResult{Answer: processor.handleReplConf(cmds)})

		case CMD_PSYNC:
			answerSlice := processor.handlePsync(cmds)

			for _, answer := range answerSlice {
				result = append(result, ProcessCmdResult{Answer: answer})
			}
		default:
			result = append(result, ProcessCmdResult{Answer: processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not able to process the cmd")})
		}
	}
	return result
}

func (processor *RespCmdProcessor) isCmdFromMaster() bool {
	if processor.replication.IsMaster() {
		return false
	}

	replicationAddress, err := replication.GetReplicationAddress(processor.connection)

	if err != nil {
		return false
	}

	return processor.replication.MasterAddress == replicationAddress
}

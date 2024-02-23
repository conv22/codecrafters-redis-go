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
	replication *replication.ReplicationStore
}

func NewRespCmdProcessor(p *resp.RespParser, storage *storage.StorageCollection, config *config.Config, replication *replication.ReplicationStore) *RespCmdProcessor {
	return &RespCmdProcessor{
		parser:      p,
		storage:     storage,
		config:      config,
		replication: replication,
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
	CMD_ACK         string = "ACK"
	CMD_GETACK      string = "GETACK"
)

type ProcessCmdResult struct {
	Answer      string
	BytesInput  []byte
	IsDuplicate bool
}

func (processor *RespCmdProcessor) getBytesInputFromCmds(cmds []resp.ParsedCmd) []byte {
	outputSlices := []resp.SliceEncoding{}

	for _, cmd := range cmds {
		outputSlices = append(outputSlices, resp.SliceEncoding{
			S: cmd.Value, Encoding: cmd.ValueType,
		})
	}

	return []byte(processor.parser.HandleEncodeSliceList(outputSlices))

}

func (processor *RespCmdProcessor) ProcessCmd(data []byte, conn net.Conn) []ProcessCmdResult {
	result := []ProcessCmdResult{}
	parsedLines, err := processor.parser.HandleParse(string(data))

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
			result = append(result, ProcessCmdResult{Answer: processor.handleSet(cmds, processor.isCmdFromMaster(conn)), IsDuplicate: true, BytesInput: processor.getBytesInputFromCmds(parsedLine)})
		case CMD_GET:
			result = append(result, ProcessCmdResult{Answer: processor.handleGet(cmds)})
		case CMD_CONFIG:
			result = append(result, ProcessCmdResult{Answer: processor.handleConfig(cmds)})
		case CMD_KEYS:
			result = append(result, ProcessCmdResult{Answer: processor.handleKeys(cmds)})
		case CMD_INFO:
			result = append(result, ProcessCmdResult{Answer: processor.handleInfo(cmds)})

		case CMD_REPLCONF:
			result = append(result, ProcessCmdResult{Answer: processor.handleReplConf(cmds, conn)})

		case CMD_PSYNC:
			answerSlice := processor.handlePsync(cmds, conn)

			for _, answer := range answerSlice {
				result = append(result, ProcessCmdResult{Answer: answer})
			}
		default:
			result = append(result, ProcessCmdResult{Answer: processor.parser.HandleEncode(RespEncodingConstants.ERROR, "not able to process the cmd")})
		}
	}
	return result
}

func (processor *RespCmdProcessor) isCmdFromMaster(conn net.Conn) bool {
	if processor.replication.IsMaster() {
		return false
	}

	replicationAddress, err := replication.GetReplicationAddress(conn)

	if err != nil {
		return false
	}

	return processor.replication.MasterAddress == replicationAddress
}

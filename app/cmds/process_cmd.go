package cmds

import (
	"github.com/codecrafters-io/redis-starter-go/app/config"
	parsers "github.com/codecrafters-io/redis-starter-go/app/parsers"
	storage "github.com/codecrafters-io/redis-starter-go/app/storage"
)

type CmdProcessor interface {
	ProcessCmd(line string) (string, error)
}

func CreateProcessor(t string, p *parsers.Parser, storage *storage.Storage, config *config.Config) CmdProcessor {
	switch t {
	case "resp":
		return &RespCmdProcessor{
			parser:  *p,
			storage: *storage,
			config:  *config,
		}

	default:
		return &RespCmdProcessor{
			parser:  *p,
			storage: *storage,
			config:  *config,
		}
	}
}

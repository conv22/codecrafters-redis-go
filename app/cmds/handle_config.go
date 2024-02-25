package cmds

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/resp"
)

const (
	configDir      = "dir"
	configFileName = "dbfilename"
)

type ConfigHandler struct {
	config *config.Config
}

func newConfigHandler(config *config.Config) *ConfigHandler {
	return &ConfigHandler{
		config: config,
	}
}

func (h *ConfigHandler) minArgs() int {
	return 2
}

func (h *ConfigHandler) processCmd(parsedResult []resp.ParsedCmd) []string {
	cmd := strings.ToLower(parsedResult[0].Value)

	switch cmd {
	case "get":
		{
			flagType := parsedResult[1].Value
			value := ""
			if flagType == configDir {
				dirFlag := h.config.DirFlag
				value = dirFlag

			}

			if flagType == configFileName {
				dbFileNameFlag := h.config.DbFilenameFlag
				value = dbFileNameFlag
			}

			encodings := []resp.SliceEncoding{
				{S: flagType, Encoding: respEncodingConstants.BULK_STRING},
				{S: value, Encoding: respEncodingConstants.BULK_STRING},
			}

			return []string{resp.HandleEncodeSliceList(encodings)}
		}
	default:
		return []string{resp.HandleEncode(respEncodingConstants.ERROR, "unsupported cmd")}
	}

}

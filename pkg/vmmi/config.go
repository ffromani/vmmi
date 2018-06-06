package vmmi

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
)

const (
	BaseConfigurationDir string = "/etc/vmmi/conf.d/"
)

func FindPluginConfigurationPath(args []string) string {
	pluginName := path.Base(args[0])
	return path.Join(BaseConfigurationDir, pluginName)
}

func (pc *PluginContext) ParseParameters(args []string) {
	if len(args) < 3 {
		details := fmt.Sprintf("expected %d arguments, received %d", 2, len(args)-1)
		pc.Abort(ErrorCodeMissingParameters, details)
	}

	pc.Params.VMid = args[1]
	pc.Params.DestinationURI = args[2]
	pc.Params.PluginConfigurationPath = FindPluginConfigurationPath(args)
	if len(args) >= 4 {
		pc.Params.PluginConfigurationPath = args[3]
	}
}

func (pc *PluginContext) ParseConfiguration() {
	var details string
	var err error
	var r io.Reader
	if pc.Params.PluginConfigurationPath == "-" {
		r = os.Stdin
	} else {
		src, err := os.Open(pc.Params.PluginConfigurationPath)
		if err != nil {
			details = fmt.Sprintf("%s", err)
			pc.Abort(ErrorCodeMalformedParameters, details)
		}
		defer src.Close()
		r = src
	}
	dec := json.NewDecoder(r)
	err = dec.Decode(pc.Config)
	if err != nil {
		details = fmt.Sprintf("%s", err)
		pc.Abort(ErrorCodeMalformedParameters, details)
	}
}

func (pc *PluginContext) Parse(args []string) {
	pc.ParseParameters(args)
	pc.ParseConfiguration()
}

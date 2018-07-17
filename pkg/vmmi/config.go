package vmmi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	BaseConfigurationDir string = "/etc/vmmi/conf.d/"
)

type URIParameters struct {
	Destination string
	Migration   string
}

type parameters struct {
	VMid                    string
	URI                     URIParameters
	PluginConfigurationPath string
}

func findPluginConfigurationPath(args []string) string {
	pluginName := path.Base(args[0])
	return path.Join(BaseConfigurationDir, pluginName)
}

func (h *Helper) parseParameters(args []string) {
	if len(args) < 4 {
		err := errors.New(fmt.Sprintf("expected at least %d arguments, received %d", 3, len(args)-1))
		h.Exit(ErrorCodeMissingParameters, err)
		return // TODO: testing helper
	}

	h.params.VMid = args[1]
	h.params.URI.Destination = args[2]
	h.params.URI.Migration = args[3]
	h.params.PluginConfigurationPath = findPluginConfigurationPath(args)
	if len(args) >= 5 {
		h.params.PluginConfigurationPath = args[4]
	}
}

func (h *Helper) readConfiguration() {
	var err error
	var r io.Reader
	if h.params.PluginConfigurationPath == "-" {
		r = os.Stdin
	} else {
		src, err := os.Open(h.params.PluginConfigurationPath)
		if err != nil {
			h.Exit(ErrorCodeBadFilePath, err)
			return // TODO: testing helper
		}
		defer src.Close()
		r = src
	}

	data, err := ioutil.ReadAll(r)
	if err != nil {
		h.Exit(ErrorCodeMalformedConfiguration, err)
		return // TODO: testing helper
	}

	h.confData = strings.NewReader(string(data))
}

func (h *Helper) parseConfiguration() {
	dec := json.NewDecoder(h.confData)
	err := dec.Decode(&h.config)
	if err != nil {
		h.Exit(ErrorCodeMalformedConfiguration, err)
		return // TODO: testing helper
	}
	// reset for next use
	h.confData.Seek(0, io.SeekStart)
}

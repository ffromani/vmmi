package vmmi

import (
	"encoding/json"
	"io"
	"os"
	"time"
)

const Version string = "0.3.0"

const (
	MessageConfiguration string = "configuration"
	MessageCompletion    string = "completion"
	MessageStatus        string = "status"
)

const (
	CompletionResultError   string = "error"
	CompletionResultSuccess string = "success"
)

type Header struct {
	Version     string `json:"vmmiVersion"`
	ContentType string `json:"contentType"`
}

type StatusMessage struct {
	Header
	Timestamp int64       `json:"timestamp"`
	Status    interface{} `json:"status"`
}

type Options struct {
	ConnectionURI string `json:"connection"`
	Verbose       int    `json:"verbose"`
}

type Parameters struct {
	VMid                    string
	DestinationURI          string
	MigrationURI            string
	PluginConfigurationPath string
}

type PluginContext struct {
	Params Parameters
	Out    io.Writer
	Config interface{}
}

func (pc *PluginContext) Status(Status interface{}) {
	msg := StatusMessage{
		Header: Header{
			Version:     Version,
			ContentType: MessageStatus,
		},
		Timestamp: time.Now().Unix(),
		Status:    &Status,
	}
	// skip errors: we have no place to report them!
	enc := json.NewEncoder(os.Stdout)
	enc.Encode(msg)
}

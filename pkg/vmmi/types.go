package vmmi

import "io"

const Version string = "0.2.0"

const MessageConfiguration string = "configuration"
const MessageCompletion string = "completion"

const CompletionResultError string = "error"
const CompletionResultSuccess string = "success"

type Header struct {
	Version     string `json:"vmmiVersion"`
	ContentType string `json:"contentType"`
}

type Options struct {
	ConnectionURI string `json:"connection"`
	Verbose       int    `json:"verbose"`
}

type Parameters struct {
	VMid                    string
	DestinationURI          string
	PluginConfigurationPath string
}

type PluginContext struct {
	Params Parameters
	Out    io.Writer
	Config interface{}
}

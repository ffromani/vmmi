package vmmi

import "io"

const Version string = "0.1.0"

const MessageConfiguration string = "configuration"
const MessageError string = "error"

type Header struct {
	Version string `json:"vmmiVersion"`
	ContentType string `json:"contentType"`
}

type Options struct {
	Connection string `json:"connection"`
	Verbose    int    `json:"verbose"`
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

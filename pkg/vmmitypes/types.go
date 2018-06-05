package vmmitypes

import "io"

const VmmiVersion string = "0.1.0"

const MessageConfiguration string = "configuration"
const MessageError string = "error"

type Header struct {
	VmmiVersion string `json:"vmmiVersion"`
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


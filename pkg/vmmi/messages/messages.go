package messages

import (
	"encoding/json"
	"io"
	"log"
	"time"
)

const Version string = "0.3.0"

const (
	ContentTypeConfiguration string = "configuration"
	ContentTypeCompletion    string = "completion"
	ContentTypeStatus        string = "status"
)

const (
	CompletionResultError   string = "error"
	CompletionResultSuccess string = "success"
)

type Header struct {
	Version     string `json:"vmmiVersion"`
	ContentType string `json:"contentType"`
}

type Credentials struct {
	Username     string `json:"username"`
	PasswordFile string `json:"passwordFile"`
}

type Options struct {
	ConnectionURI         string      `json:"connection"`
	Verbose               int         `json:"verbose"`
	LogFilePath           string      `json:"logFilePath"`
	ConnectionCredentials Credentials `json:"connectionCredentials"`
}

type Configuration struct {
	Header
	Configuration Options `json:"configuration"`
}

type Status struct {
	Header
	Timestamp int64       `json:"timestamp"`
	Status    interface{} `json:"status"`
}

type SuccessData struct {
}

type CompletionSuccessData struct {
	Result  string       `json:"result"`
	Success *SuccessData `json:"success"`
}

type CompletionSuccess struct {
	Header
	Timestamp  int64                  `json:"timestamp"`
	Completion *CompletionSuccessData `json:"completion"`
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type CompletionErrorData struct {
	Result string     `json:"result"`
	Error  *ErrorData `json:"error"`
}

type CompletionError struct {
	Header
	Timestamp  int64                `json:"timestamp"`
	Completion *CompletionErrorData `json:"completion"`
}

type Sink struct {
	W io.Writer
	L *log.Logger
}

func (s *Sink) ReportSuccess() {
	msg := CompletionSuccess{
		Header: Header{
			Version:     Version,
			ContentType: ContentTypeCompletion,
		},
		Timestamp: time.Now().Unix(),
		Completion: &CompletionSuccessData{
			Result:  CompletionResultSuccess,
			Success: &SuccessData{},
		},
	}
	// TODO
	// skip errors: we have no place to report them!
	enc := json.NewEncoder(s.W)
	enc.Encode(msg)
}

func (s *Sink) ReportError(code int, message string, details string) {
	msg := CompletionError{
		Header: Header{
			Version:     Version,
			ContentType: ContentTypeCompletion,
		},
		Timestamp: time.Now().Unix(),
		Completion: &CompletionErrorData{
			Result: CompletionResultError,
			Error: &ErrorData{
				Code:    code,
				Message: message,
				Details: details,
			},
		},
	}
	// skip errors: we have no place to report them!
	enc := json.NewEncoder(s.W)
	enc.Encode(msg)
}

func (s *Sink) ReportStatus(status interface{}) {
	msg := Status{
		Header: Header{
			Version:     Version,
			ContentType: ContentTypeStatus,
		},
		Status:    status,
		Timestamp: time.Now().Unix(),
	}
	// TODO handle error
	enc := json.NewEncoder(s.W)
	enc.Encode(msg)
}

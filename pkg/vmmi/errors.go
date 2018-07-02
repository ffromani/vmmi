package vmmi

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	ErrorCodeNone = iota
	ErrorCodeUnknown
	ErrorCodeOperationFailed
	ErrorCodeBadFilePath
	ErrorCodeMalformedParameters
	ErrorCodeMalformedConfiguration
	ErrorCodeMissingParameters
	ErrorCodeMigrationFailed
	ErrorCodeMigrationAborted
	ErrorCodeVMUnknown
	ErrorCodeVMDisappeared
	ErrorCodeLibvirtDisconnected
)

type SuccessData struct {
}

type SuccessCompletionData struct {
	Result  string       `json:"result"`
	Success *SuccessData `json:"success"`
}

type SuccessCompletionMessage struct {
	Header
	Timestamp  int64                  `json:"timestamp"`
	Completion *SuccessCompletionData `json:"completion"`
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type ErrorCompletionData struct {
	Result string     `json:"result"`
	Error  *ErrorData `json:"error"`
}

type ErrorCompletionMessage struct {
	Header
	Timestamp  int64                `json:"timestamp"`
	Completion *ErrorCompletionData `json:"completion"`
}

func Strerror(code int) string {
	switch code {
	case ErrorCodeNone:
		return ""
	case ErrorCodeUnknown:
		return "unexpected error"
	case ErrorCodeOperationFailed:
		return "operation failed"
	case ErrorCodeBadFilePath:
		return "bad file path specification"
	case ErrorCodeMalformedParameters:
		return "malformed parameters"
	case ErrorCodeMalformedConfiguration:
		return "malformed configuration"
	case ErrorCodeMissingParameters:
		return "missing parameters"
	case ErrorCodeMigrationFailed:
		return "libvirt migration failed"
	case ErrorCodeMigrationAborted:
		return "migration aborted"
	case ErrorCodeVMUnknown:
		return "VM unknown"
	case ErrorCodeVMDisappeared:
		return "VM disappeared"
	case ErrorCodeLibvirtDisconnected:
		return "Lost connection to libvirt"
	}
	return "unknown"
}

func ReportSuccess(w io.Writer) {
	msg := SuccessCompletionMessage{
		Header: Header{
			Version:     Version,
			ContentType: MessageCompletion,
		},
		Timestamp: time.Now().Unix(),
		Completion: &SuccessCompletionData{
			Result:  CompletionResultSuccess,
			Success: &SuccessData{},
		},
	}
	// skip errors: we have no place to report them!
	enc := json.NewEncoder(w)
	enc.Encode(msg)
}

func ReportError(w io.Writer, code int, details string) {
	msg := ErrorCompletionMessage{
		Header: Header{
			Version:     Version,
			ContentType: MessageCompletion,
		},
		Timestamp: time.Now().Unix(),
		Completion: &ErrorCompletionData{
			Result: CompletionResultError,
			Error: &ErrorData{
				Code:    code,
				Message: Strerror(code),
				Details: details,
			},
		},
	}
	// skip errors: we have no place to report them!
	enc := json.NewEncoder(w)
	enc.Encode(msg)
}

func (pc *PluginContext) CompleteWithErrorDetails(code int, details string) {
	ReportError(pc.Out, code, details)
	os.Exit(1)
}

func (pc *PluginContext) CompleteWithErrorValue(code int, err error) {
	details := fmt.Sprintf("%s", err)
	pc.CompleteWithErrorDetails(code, details)
}

func (pc *PluginContext) CompleteWithSuccess() {
	ReportSuccess(pc.Out)
	os.Exit(0)
}

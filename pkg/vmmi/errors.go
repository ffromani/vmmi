package vmmi

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

const (
	ErrorCodeNone = iota
	ErrorCodeUnknown
	ErrorCodeOperationFailed
	ErrorCodeMalformedParameters
	ErrorCodeMissingParameters
	ErrorCodeMigrationFailed
	ErrorCodeMigrationAborted
	ErrorCodeVMUnknown
	ErrorCodeVMDisappeared
	ErrorCodeLibvirtDisconnected
)

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type ErrorMessage struct {
	Header
	Timestamp int64      `json:"timestamp"`
	Error     *ErrorData `json:"error"`
}

func Strerror(code int) string {
	switch code {
	case ErrorCodeNone:
		return ""
	case ErrorCodeUnknown:
		return "unexpected error"
	case ErrorCodeOperationFailed:
		return "operation failed"
	case ErrorCodeMalformedParameters:
		return "malformed parameters"
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

func Report(w io.Writer, code int, details string) {
	msg := ErrorMessage{
		Header: Header{
			Version: Version,
			ContentType: MessageError,
		},
		Timestamp: time.Now().Unix(),
		Error: &ErrorData{
			Code:    code,
			Message: Strerror(code),
			Details: details,
		},
	}
	// skip errors: we have no place to report them!
	enc := json.NewEncoder(w)
	enc.Encode(msg)
}

func (pc *PluginContext) Report(code int, details string) {
	Report(pc.Out, code, details)
}


func (pc *PluginContext) ReportError(code int, err error) {
	details := fmt.Sprintf("%s", err)
	pc.Report(code, details)
}

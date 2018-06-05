package vmmierrors

import (
	"encoding/json"
	"github.com/fromanirh/vmmi/pkg/vmmitypes"
	"io"
	"os"
	"time"
)

const (
	ErrorCodeNone = iota
	ErrorCodeMalformedParameters
	ErrorCodeMissingParameters
	ErrorCodeMigrationFailed
	ErrorCodeMigrationAborted
	ErrorCodeVMDisappeared
)

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type ErrorMessage struct {
	vmmitypes.Header
	Timestamp int64      `json:"timestamp"`
	Error     *ErrorData `json:"error"`
}

func Strerror(code int) string {
	switch code {
	case ErrorCodeNone:
		return "none"
	case ErrorCodeMalformedParameters:
		return "malformed parameters"
	case ErrorCodeMissingParameters:
		return "missing parameters"
	case ErrorCodeMigrationFailed:
		return "libvirt migration failed"
	case ErrorCodeMigrationAborted:
		return "migration aborted"
	case ErrorCodeVMDisappeared:
		return "VM disappeared"
	}
	return "unknown"
}

func Report(w io.Writer, code int, details string) {
	msg := ErrorMessage{
		Header: vmmitypes.Header{
			VmmiVersion: vmmitypes.VmmiVersion,
			ContentType: vmmitypes.MessageError,
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

func Abort(pc *vmmitypes.PluginContext, code int, details string) {
	Report(pc.Out, code, details)
	os.Exit(1)
}

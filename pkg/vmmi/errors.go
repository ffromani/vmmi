package vmmi

const (
	ErrorCodeNone = iota
	ErrorCodeUnknown
	ErrorCodeOperationFailed
	ErrorCodeBadFilePath
	ErrorCodeConfigurationFailed
	ErrorCodeMalformedParameters
	ErrorCodeMalformedConfiguration
	ErrorCodeMissingParameters
	ErrorCodeMigrationFailed
	ErrorCodeMigrationAborted
	ErrorCodeVMUnknown
	ErrorCodeVMDisappeared
	ErrorCodeLibvirtDisconnected
)

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
	case ErrorCodeConfigurationFailed:
		return "failed to apply the configuration"
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

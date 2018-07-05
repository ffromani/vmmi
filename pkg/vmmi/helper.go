package vmmi

import (
	"fmt"
	"github.com/fromanirh/vmmi/pkg/vmmi/messages"
	libvirt "github.com/libvirt/libvirt-go"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Helper struct {
	logFile   *os.File
	logHandle *log.Logger
	confData  *strings.Reader
	config    messages.Configuration
	params    parameters
	dom       *libvirt.Domain
	conn      *libvirt.Connect
}

func (h *Helper) Log() *log.Logger {
	return h.logHandle
}

func NewHelper(args []string) *Helper {
	var logsink io.Writer = ioutil.Discard
	var err error

	h := Helper{}
	h.parseParameters(args)
	h.readConfiguration()
	h.parseConfiguration()

	if h.config.Configuration.LogFilePath != "" {
		h.logFile, err = os.OpenFile(h.config.Configuration.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		logsink = h.logFile
	}

	if err != nil {
		h.completeWithErrorValue(ErrorCodeBadFilePath, err)
	}
	h.logHandle = log.New(logsink, "touri: ", log.LstdFlags)

	h.Log().Printf("args: %#v", args)
	h.Log().Printf("configuration: %#v", h.config)

	h.conn, err = connect(&h.config.Configuration)
	if err != nil {
		h.completeWithErrorValue(ErrorCodeLibvirtDisconnected, err)
	}

	h.Log().Printf("connected")

	h.dom, err = h.conn.LookupDomainByUUIDString(h.params.VMid)
	if err != nil {
		h.completeWithErrorValue(ErrorCodeVMUnknown, err)
	}
	h.Log().Printf("lookup succesfull")
	return &h
}

func (h *Helper) Close() error {
	if h.dom != nil {
		h.dom.Free()
		h.dom = nil
	}
	if h.conn != nil {
		h.conn.Close()
		h.conn = nil
	}
	h.Log().Printf("done!")
	if h.logFile != nil {
		h.logFile.Close()
		h.logFile = nil
	}
	return nil
}

func (h *Helper) GetDomain() *libvirt.Domain {
	return h.dom
}

func (h *Helper) GetURIParameters() URIParameters {
	return h.params.URI
}

func (h *Helper) sendStatus(mon Monitor) {
	payload, err := mon.Status()
	if err != nil {
		// TODO report error
		return
	}

	sink := messages.Sink{W: os.Stdout, L: h.Log()}
	sink.ReportStatus(payload)
}

func (h *Helper) completeWithErrorDetails(code int, details string) {
	sink := messages.Sink{W: os.Stderr, L: h.Log()}
	sink.ReportError(code, Strerror(code), details)
	os.Exit(1)
}

func (h *Helper) completeWithErrorValue(code int, err error) {
	details := fmt.Sprintf("%s", err)
	h.completeWithErrorDetails(code, details)
}

func (h *Helper) completeWithSuccess() {
	sink := messages.Sink{W: os.Stderr, L: h.Log()}
	sink.ReportSuccess()
	os.Exit(0)
}

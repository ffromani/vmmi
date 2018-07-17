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
	exitFunc  func(code int)
	outsink   io.Writer
	errsink   io.Writer
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

func newHelper() *Helper {
	h := &Helper{
		exitFunc: os.Exit,
		outsink:  os.Stdout,
		errsink:  os.Stderr,
	}
	return h
}

func NewHelper(args []string) *Helper {
	h := newHelper()
	h.parseParameters(args)
	h.readConfiguration()
	h.parseConfiguration()
	h.openLog()

	h.Log().Printf("args: %#v", args)
	h.Log().Printf("configuration: %#v", h.config)

	h.connectToLibvirt()

	return h
}

func (h *Helper) openLog() *Helper {
	var logsink io.Writer = ioutil.Discard
	var err error

	if h.config.Configuration.LogFilePath != "" {
		h.logFile, err = os.OpenFile(h.config.Configuration.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		logsink = h.logFile
	}

	if err != nil {
		h.Exit(ErrorCodeBadFilePath, err)
	}
	h.logHandle = log.New(logsink, "touri: ", log.LstdFlags)
	return h
}

func (h *Helper) connectToLibvirt() *Helper {
	var err error

	h.conn, err = connect(&h.config.Configuration)
	if err != nil {
		h.Exit(ErrorCodeLibvirtDisconnected, err)
	}

	h.Log().Printf("connected")

	h.dom, err = h.conn.LookupDomainByUUIDString(h.params.VMid)
	if err != nil {
		h.Exit(ErrorCodeVMUnknown, err)
	}
	h.Log().Printf("lookup succesfull")

	return h
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

func (h *Helper) Domain() *libvirt.Domain {
	return h.dom
}

func (h *Helper) URIParameters() URIParameters {
	return h.params.URI
}

func (h *Helper) sendStatus(mon Monitor) {
	msg, err := mon.Status(messages.NewStatus())
	if err != nil {
		// TODO report error
		return
	}

	sink := messages.Sink{W: h.outsink, L: h.Log()}
	sink.ReportStatus(msg)
}

func (h *Helper) Exit(code int, err ...error) {
	sink := messages.Sink{W: h.errsink, L: h.Log()}
	if code == ErrorCodeNone {
		sink.ReportSuccess()
		h.ExitFunc(0)
	}
	sink.ReportError(code, Strerror(code), fmt.Sprintf("%s", err[0]))
	h.ExitFunc(1)
}

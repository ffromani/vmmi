package main

import (
	"fmt"
	"github.com/fromanirh/vmmi/pkg/vmmi"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Options struct {
	vmmi.Options
	Delay string `json:"delay"`
}

type PluginConfiguration struct {
	vmmi.Header
	Configuration Options `json:"configuration"`
}

func main() {
	var details string
	conf := PluginConfiguration{}
	pc := &vmmi.PluginContext{
		Config: &conf,
		Out:    os.Stdout,
	}
	pc.Parse(os.Args)

	delay, err := time.ParseDuration(conf.Configuration.Delay)
	if err != nil {
		details = fmt.Sprintf("bad delay specification: %s", conf.Configuration.Delay)
		pc.Abort(vmmi.ErrorCodeMalformedParameters, details)
	}

	errCode := vmmi.ErrorCodeMigrationFailed

	t := time.NewTimer(delay)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTERM)

	start := time.Now()
	select {
	case s := <-sigs:
		switch s {
		case syscall.SIGINT, syscall.SIGSTOP:
			errCode = vmmi.ErrorCodeMigrationAborted
		case syscall.SIGTERM:
			errCode = vmmi.ErrorCodeNone
		default:
			errCode = vmmi.ErrorCodeUnknown
		}
	case <-t.C:
		// do nothing
	}
	stop := time.Now()

	details = fmt.Sprintf("cannot migrate VM %s to %s using %s (took %v)", pc.Params.VMid, pc.Params.DestinationURI, conf.Configuration.Connection, stop.Sub(start))
	pc.Abort(errCode, details)
}

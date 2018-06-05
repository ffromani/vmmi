package main

import (
	"fmt"
	"github.com/fromanirh/vmmi/pkg/vmmiconfig"
	"github.com/fromanirh/vmmi/pkg/vmmierrors"
	"github.com/fromanirh/vmmi/pkg/vmmitypes"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Options struct {
	vmmitypes.Options
	Delay string `json:"delay"`
}

type PluginConfiguration struct {
	vmmitypes.Header
	Configuration Options `json:"configuration"`
}

// TODO: handle signals

func main() {
	var details string
	conf := PluginConfiguration{}
	pc := &vmmitypes.PluginContext{
		Config: &conf,
		Out:    os.Stdout,
	}
	vmmiconfig.Parse(pc, os.Args)

	delay, err := time.ParseDuration(conf.Configuration.Delay)
	if err != nil {
		details = fmt.Sprintf("bad delay specification: %s", conf.Configuration.Delay)
		vmmierrors.Abort(pc, vmmierrors.ErrorCodeMalformedParameters, details)
	}

	errCode := vmmierrors.ErrorCodeMigrationFailed

	t := time.NewTimer(delay)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTERM)

	start := time.Now()
	select {
	case s := <-sigs:
		switch s {
		case syscall.SIGINT, syscall.SIGSTOP:
			errCode = vmmierrors.ErrorCodeMigrationAborted
		case syscall.SIGTERM:
			errCode = vmmierrors.ErrorCodeNone
		default:
			errCode = vmmierrors.ErrorCodeUnknown
		}
	case <- t.C:
		// do nothing
	}
	stop := time.Now()

	details = fmt.Sprintf("cannot migrate VM %s to %s using %s (took %v)", pc.Params.VMid, pc.Params.DestinationURI, conf.Configuration.Connection, stop.Sub(start))
	vmmierrors.Abort(pc, errCode, details)
}

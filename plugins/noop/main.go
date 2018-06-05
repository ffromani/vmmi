package main

import (
	"fmt"
	"github.com/fromanirh/vmmi/pkg/vmmiconfig"
	"github.com/fromanirh/vmmi/pkg/vmmierrors"
	"github.com/fromanirh/vmmi/pkg/vmmitypes"
	"os"
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

	time.Sleep(delay)

	details = fmt.Sprintf("cannot migrate VM %s to %s using %s (took %v)", pc.Params.VMid, pc.Params.DestinationURI, conf.Configuration.Connection, conf.Configuration.Delay)
	vmmierrors.Abort(pc, vmmierrors.ErrorCodeMigrationFailed, details)
}

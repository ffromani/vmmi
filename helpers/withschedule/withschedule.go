package main

import (
	"github.com/fromanirh/vmmi/pkg/convsched"
	"github.com/fromanirh/vmmi/pkg/vmmi"
	libvirt "github.com/libvirt/libvirt-go"
	"os"
)

type SchedulingMigrator struct {
	MigrationURI   string
	DestinationURI string
	Domain         *libvirt.Domain
	monitor        *convsched.SchedulingMonitor
}

func (sm *SchedulingMigrator) Run(resChan chan error) {
	params := libvirt.DomainMigrateParameters{
		URI:    sm.MigrationURI,
		URISet: true,
	}
	flags := libvirt.MIGRATE_LIVE | libvirt.MIGRATE_PEER2PEER | libvirt.MIGRATE_PERSIST_DEST | libvirt.MIGRATE_COMPRESSED | libvirt.MIGRATE_AUTO_CONVERGE
	if sm.monitor.ScheduleHasPostcopy() {
		flags |= libvirt.MIGRATE_POSTCOPY
	}
	resChan <- sm.Domain.MigrateToURI3(sm.DestinationURI, &params, flags)
}

func main() {
	h := vmmi.NewHelper(os.Args)
	uris := h.URIParameters()
	mon := convsched.SchedulingMonitor{
		Domain:         h.Domain(),
		DestinationURI: uris.Destination,
		MigrationURI:   uris.Migration,
		Log:            h.Log(),
	}
	mig := SchedulingMigrator{
		Domain:         h.Domain(),
		DestinationURI: uris.Destination,
		MigrationURI:   uris.Migration,
		monitor:        &mon,
	}
	h.WaitForCompletion(&mon, &mig)
}

package main

import (
	"github.com/fromanirh/vmmi/pkg/vmmi"
	libvirt "github.com/libvirt/libvirt-go"
	"io"
	"os"
)

type NOPMonitor struct {
}

func (n *NOPMonitor) Configure(r io.Reader) error {
	return nil
}

func (n *NOPMonitor) Run(resChan chan error) {
	resChan <- nil
}

func (n *NOPMonitor) Stop() {
}

func (n *NOPMonitor) Status() (interface{}, error) {
	return nil, nil
}

type TrivialMigrator struct {
	MigrationURI   string
	DestinationURI string
	Domain         *libvirt.Domain
}

func (tm *TrivialMigrator) Run(resChan chan error) {
	params := libvirt.DomainMigrateParameters{
		URI:    tm.MigrationURI,
		URISet: true,
	}
	flags := libvirt.MIGRATE_LIVE | libvirt.MIGRATE_PEER2PEER
	resChan <- tm.Domain.MigrateToURI3(tm.DestinationURI, &params, flags)
}

func main() {
	h := vmmi.NewHelper(os.Args)
	mon := NOPMonitor{}
	uris := h.URIParameters()
	mig := TrivialMigrator{
		Domain:         h.Domain(),
		DestinationURI: uris.Destination,
		MigrationURI:   uris.Migration,
	}
	h.WaitForCompletion(&mon, &mig)
}

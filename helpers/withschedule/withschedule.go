package main

import (
	"encoding/json"
	"github.com/fromanirh/vmmi/pkg/vmmi"
	"github.com/fromanirh/vmmi/pkg/vmmi/messages"
	libvirt "github.com/libvirt/libvirt-go"
	"io"
	"log"
	"os"
)

type ConvergenceAction struct {
	Name   string   `json:"name"`
	Params []string `json:"params"`
}

type ConvergenceItem struct {
	Action ConvergenceAction `json:"action"`
	Limit  int               `json:"limit"`
}

type ConvergenceSchedule struct {
	Init     ConvergenceAction   `json:"init"`
	Stalling []ConvergenceAction `json:"stalling"`
}

type SchedulingMonitor struct {
	lh    *log.Logger
	sched ConvergenceSchedule
}

func (mon *SchedulingMonitor) Configure(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&mon.sched)
}

func (mon *SchedulingMonitor) Run(resChan chan error) {
	mon.lh.Printf("%#v", mon.sched)
	resChan <- nil
}

func (mon *SchedulingMonitor) Stop() {
}

func (mon *SchedulingMonitor) Status(msg *messages.Status) (interface{}, error) {
	return msg, nil
}

func (mon *SchedulingMonitor) ScheduleHasPostcopy() bool {
	return true
}

type SchedulingMigrator struct {
	MigrationURI   string
	DestinationURI string
	Domain         *libvirt.Domain
	monitor        *SchedulingMonitor
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
	mon := SchedulingMonitor{
		lh: h.Log(),
	}
	uris := h.URIParameters()
	mig := SchedulingMigrator{
		Domain:         h.Domain(),
		DestinationURI: uris.Destination,
		MigrationURI:   uris.Migration,
		monitor:        &mon,
	}
	h.WaitForCompletion(&mon, &mig)
}

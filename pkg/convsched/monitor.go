package convsched

import (
	"github.com/fromanirh/vmmi/pkg/vmmi/messages"
	"github.com/fromanirh/vmmi/pkg/vmmi/progress"
	libvirt "github.com/libvirt/libvirt-go"
	"io"
	"log"
	"time"
)

const (
	PostCopyPhaseNone = iota
	PostCopyPhaseRequested
	PostCopyPhaseRunning
)

type SchedulingMonitor struct {
	MigrationURI   string
	DestinationURI string
	Domain         *libvirt.Domain
	Log            *log.Logger
	schedule       *ConvergenceSchedule
	interval       time.Duration
	stopped        chan bool
}

type domainMigrator struct {
	mon *SchedulingMonitor
}

func (d *domainMigrator) Progress() *progress.Progress {
	return progress.NewProgress(d.mon.Domain)
}

func (d *domainMigrator) SetDowntime(value int) error {
	// TODO: downtime must be >= 0
	d.mon.Log.Printf("action: setting downtime to %v", value)
	return d.mon.Domain.MigrateSetMaxDowntime(uint64(value), 0)
}

func (d *domainMigrator) StartPostCopy() error {
	d.mon.Log.Printf("action: switching to post copy")
	return d.mon.Domain.MigrateStartPostCopy(0)
}

func (d *domainMigrator) Abort() error {
	d.mon.Log.Printf("action: aborting migration")
	err := d.mon.Domain.AbortJob()
	d.mon.Stop()
	return err
}

func (mon *SchedulingMonitor) Configure(r io.Reader) error {
	mon.stopped = make(chan bool, 1) // TODO move into Run() ?
	conf, err := LoadConfiguration(r)
	if err == nil {
		mon.schedule = &conf.Schedule
		mon.interval = conf.MonitorInterval
	}
	return err
}

func (mon *SchedulingMonitor) Stop() {
	mon.stopped <- true
}

func (mon *SchedulingMonitor) Status(msg *messages.Status) (interface{}, error) {
	return msg, nil
}

func (mon *SchedulingMonitor) ScheduleHasPostcopy() bool {
	return mon.schedule.HasPostcopy()
}

func (mon *SchedulingMonitor) Run(resChan chan error) {
	mig := &domainMigrator{mon: mon}
	err := mon.executeInit(mig)
	if err == nil {
		err = mon.runLoop(mig)
	}
	resChan <- err
}

type monitorInfo struct {
	postCopyPhase     int
	lowMark           int64
	lastDataRemaining int64
	iterationCount    int64
	step              uint64
}

func (mon *SchedulingMonitor) runLoop(mig VMMigrator) error {
	var err error
	monInfo := monitorInfo{
		postCopyPhase:     PostCopyPhaseNone,
		lowMark:           -1,
		lastDataRemaining: -1,
	}
	ticker := time.NewTicker(mon.interval * time.Second)
	stopped := false
	for !stopped {
		select {
		case stopped = <-mon.stopped:
			// nothing to do there
		case <-ticker.C:
			err = mon.runStep(&monInfo, mig)
			if err != nil {
				stopped = true
			}
			monInfo.step++
		}
	}
	return err
}

func (mon *SchedulingMonitor) runStep(monInfo *monitorInfo, mig VMMigrator) error {
	var err error
	prog := mig.Progress()
	if prog == nil {
		// not ready yet; not critical, let's try again later
		return err
	}
	dataRemaining := prog.DataRemaining()

	mon.Log.Printf("step %#v with data remaining %v", monInfo, dataRemaining)

	if monInfo.postCopyPhase != PostCopyPhaseNone {
		if monInfo.postCopyPhase == PostCopyPhaseRunning {
			mon.Log.Printf("Post-copy migration still in progress: %v", dataRemaining)
		}
	} else if monInfo.lowMark == -1 || monInfo.lowMark > dataRemaining {
		monInfo.lowMark = dataRemaining
	} else {
		mon.Log.Printf("Migration stalling: remaining (%vMiB) > lowmark (%vMiB).",
			dataRemaining/1024./1024.,
			monInfo.lowMark/1024./1024.)

	}

	if monInfo.postCopyPhase == PostCopyPhaseNone && monInfo.lastDataRemaining != -1 && monInfo.lastDataRemaining < dataRemaining {
		monInfo.iterationCount++
		mon.Log.Printf("New iteration detected: %v, remaining convergence schedule %v", monInfo.iterationCount, mon.schedule)

		action := mon.schedule.PopAction(monInfo.iterationCount)
		if action != nil {
			mon.Log.Printf("loop: applying convergence action '%v'", action)
			err = action.Exec(mig)
		}
	}

	monInfo.lastDataRemaining = dataRemaining
	return err
}

func (mon *SchedulingMonitor) executeInit(mig VMMigrator) error {
	var err error
	for _, action := range mon.schedule.Init {
		mon.Log.Printf("init: applying convergence action '%v'", action)
		err = action.Exec(mig)
		if err != nil {
			return err
		}
	}
	return err
}

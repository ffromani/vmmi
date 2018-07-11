package convsched

import (
	"github.com/fromanirh/vmmi/pkg/vmmi/messages"
	"github.com/fromanirh/vmmi/pkg/vmmi/progress"
	libvirt "github.com/libvirt/libvirt-go"
	"io"
	"log"
	"strconv"
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
	err := mon.executeInit()
	if err == nil {
		err = mon.runLoop()
	}
	resChan <- err
}

type monitorInfo struct {
	postCopyPhase     int
	lowMark           int64
	lastDataRemaining int64
	iterationCount    int64
}

func (mon *SchedulingMonitor) runLoop() error {
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
			err = mon.runStep(&monInfo)
			if err != nil {
				stopped = true
			}
		}
	}
	return err
}

func (mon *SchedulingMonitor) runStep(monInfo *monitorInfo) error {
	prog := progress.NewProgress(mon.Domain)
	if prog == nil {
		// not ready yet; not critical, let's try again later
		return nil
	}
	info := prog.JobInfo()
	if info == nil {
		return nil
	}
	dataRemaining := int64(info.DataRemaining)

	if monInfo.postCopyPhase != PostCopyPhaseNone {
		if monInfo.postCopyPhase == PostCopyPhaseRunning {
			mon.Log.Printf("Post-copy migration still in progress: %d", info.DataRemaining)
		}
	} else if monInfo.lowMark == -1 || monInfo.lowMark > dataRemaining {
		monInfo.lowMark = dataRemaining
	} else {
		mon.Log.Printf("Migration stalling: remaining (%vMiB) > lowmark (%vMiB).",
			info.DataRemaining/1024./1024.,
			monInfo.lowMark/1024./1024.)

	}

	if monInfo.postCopyPhase == PostCopyPhaseNone && monInfo.lastDataRemaining != -1 && monInfo.lastDataRemaining < dataRemaining {
		monInfo.iterationCount++
		mon.Log.Printf("New iteration detected: %v", monInfo.iterationCount)
		mon.executeActionForIteration(monInfo.iterationCount)
	}

	monInfo.lastDataRemaining = int64(info.DataRemaining)
	mon.Log.Printf("progress: %v", prog)
	return nil
}

func (mon *SchedulingMonitor) executeInit() error {
	var err error
	for _, action := range mon.schedule.Init {
		err = mon.executeAction(action)
		if err != nil {
			return err
		}
	}
	return err
}

func (mon *SchedulingMonitor) executeActionForIteration(currentIteration int64) error {
	var err error
	head := mon.schedule.Stalling[0]

	mon.Log.Printf("Stalling for %v iterations, checking to make next action: %v", currentIteration, head)
	if head.Limit < currentIteration {
		err = mon.executeAction(head.Action)
		mon.schedule.Stalling = mon.schedule.Stalling[1:]
		mon.Log.Printf("setting convergence schedule to: %v", mon.schedule.Stalling)
	}
	return err
}

func (mon *SchedulingMonitor) executeAction(action ConvergenceAction) error {
	var err error
	switch action.Name {
	case ActionSetDowntime:
		downtime, err := strconv.Atoi(action.Params[0])
		if err != nil {
			return err
		}
		mon.Log.Printf("action: setting downtime to %v", downtime)
		err = mon.Domain.MigrateSetMaxDowntime(uint64(downtime), 0)
	case ActionEnablePostCopy:
		mon.Log.Printf("action: switching to post copy")
		err = mon.Domain.MigrateStartPostCopy(0)
	case ActionAbort:
		mon.Log.Printf("action: aborting migration")
		err = mon.Domain.AbortJob()
		mon.Stop()
	}
	return err
}

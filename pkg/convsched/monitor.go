package convsched

import (
	"github.com/fromanirh/vmmi/pkg/vmmi/messages"
	"io"
)

type SchedulingMonitor struct {
	sched *ConvergenceSchedule
}

func (mon *SchedulingMonitor) Configure(r io.Reader) error {
	cs, err := Load(r)
	if err == nil {
		mon.sched = cs
	}
	return err
}

func (mon *SchedulingMonitor) Run(resChan chan error) {
	resChan <- nil
}

func (mon *SchedulingMonitor) Stop() {
}

func (mon *SchedulingMonitor) Status(msg *messages.Status) (interface{}, error) {
	return msg, nil
}

func (mon *SchedulingMonitor) ScheduleHasPostcopy() bool {
	return mon.sched.HasPostcopy()
}

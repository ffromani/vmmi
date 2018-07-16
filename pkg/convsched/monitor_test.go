package convsched

import (
	"errors"
	"github.com/fromanirh/vmmi/pkg/vmmi/progress"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestScheduleLoadConfiguration(t *testing.T) {
	confFile := "conf_example_simple.json"
	f, err := os.Open(confFile)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}
	defer f.Close()

	mon := SchedulingMonitor{}
	err = mon.Configure(f)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}

	if mon.schedule == nil || len(mon.schedule.Init) == 0 || len(mon.schedule.Stalling) == 0 {
		t.Errorf("sched: %#v", mon.schedule)
	}
}

func TestScheduleHasNotPostcopy(t *testing.T) {
	confFile := "conf_example_simple.json"
	f, err := os.Open(confFile)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}
	defer f.Close()

	mon := SchedulingMonitor{}
	err = mon.Configure(f)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}

	if mon.ScheduleHasPostcopy() {
		t.Errorf("%v has postcopy, but it should not", confFile)
	}
}

func TestScheduleHasPostcopy(t *testing.T) {
	confFile := "conf_example_simple_postcopy.json"
	f, err := os.Open(confFile)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}
	defer f.Close()

	mon := SchedulingMonitor{}
	err = mon.Configure(f)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}

	if !mon.ScheduleHasPostcopy() {
		t.Errorf("%v has not postcopy, but it should", confFile)
	}
}

type fakeMigrator struct {
	dataRemaining []int64 // amount of remaining data. Will decrease one unit each step
	mon           *SchedulingMonitor
	Downtime      int
	Postcopy      bool
	Aborted       bool
}

func (f *fakeMigrator) Progress() *progress.Progress {
	value := f.dataRemaining[0]
	f.dataRemaining = f.dataRemaining[1:]
	return &progress.Progress{UserDataRemaining: value}
}

func (f *fakeMigrator) SetDowntime(value int) error {
	f.Downtime = value
	return nil
}

func (f *fakeMigrator) StartPostCopy() error {
	f.Postcopy = true
	return nil
}

func (f *fakeMigrator) Abort() error {
	f.Aborted = true
	if f.mon != nil {
		f.mon.Stop()
	}
	return nil
}

type failingMigrator struct {
}

func (f *failingMigrator) Progress() *progress.Progress {
	return nil
}

func (f *failingMigrator) SetDowntime(value int) error {
	return errors.New("SetDowntime failed")
}

func (f *failingMigrator) StartPostCopy() error {
	return errors.New("StartPostCopy failed")
}

func (f *failingMigrator) Abort() error {
	return errors.New("Abort failed")
}

func TestScheduleExecuteInit(t *testing.T) {
	confFile := "conf_abort_init.json"
	f, err := os.Open(confFile)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}
	defer f.Close()

	mon := SchedulingMonitor{
		Log: log.New(os.Stderr, "test: ", log.LstdFlags),
	}
	err = mon.Configure(f)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}

	mig := &fakeMigrator{}
	err = mon.executeInit(mig)
	if err != nil {
		t.Errorf("executeInit failed: %v", err)
	}
	if !mig.Aborted {
		t.Errorf("schedule not aborted: %v", mon.schedule)
	}

	fail := &failingMigrator{}
	err = mon.executeInit(fail)
	if err == nil {
		t.Errorf("executeInit should have failed, but it did'nt: %v", err)
	}
}

func TestRunLoopJustOnce(t *testing.T) {
	confFile := "conf_abort_init.json"
	f, err := os.Open(confFile)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}
	defer f.Close()

	mon := SchedulingMonitor{
		Log: log.New(ioutil.Discard, "test: ", log.LstdFlags),
	}
	err = mon.Configure(f)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}

	mig := &fakeMigrator{
		dataRemaining: []int64{10, 12, 9, 8, 7, 6, 5},
		mon:           &mon,
	}
	err = mon.runLoop(mig)
	if err != nil {
		t.Errorf("runLoop failed: %v", err)
	}
	if !mig.Aborted {
		t.Errorf("schedule not aborted: %v", mon.schedule)
	}
}

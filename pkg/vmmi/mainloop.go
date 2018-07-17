package vmmi

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (h *Helper) WaitForCompletion(mon Monitor, mig Migrator) {
	defer h.Close()
	var err error

	err = mon.Configure(h.confData)
	if err != nil {
		h.Exit(ErrorCodeConfigurationFailed, err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGSTOP, syscall.SIGTERM, syscall.SIGUSR1)

	h.Log().Printf("signals set")

	migrationError := make(chan error, 1)
	go mig.Run(migrationError)
	h.Log().Printf("migration started")

	monitorError := make(chan error, 1)
	go mon.Run(monitorError)
	h.Log().Printf("monitor started")

	errCode := ErrorCodeNone
	start := time.Now()

	h.Log().Printf("waiting")
	done := false
	for !done {
		select {
		case s := <-sigs:
			err := errors.New(fmt.Sprintf("interrupted by signal %v", s))

			switch s {
			case syscall.SIGINT, syscall.SIGSTOP:
				mon.Stop()
				err = h.dom.AbortJob()
				if err != nil {
					h.Exit(ErrorCodeOperationFailed, err)
				}
				errCode = ErrorCodeMigrationAborted
				done = true
			case syscall.SIGTERM:
				mon.Stop()
				errCode = ErrorCodeNone
				done = true
			case syscall.SIGUSR1:
				h.sendStatus(mon)
			}

		case err = <-migrationError:
			done = true
			mon.Stop() // ensure monitor is stopped
			h.Log().Printf("migration stop err=%v", err)
			if err == nil {
				errCode = ErrorCodeNone
				h.Log().Printf("migration completed in %v", time.Now().Sub(start))
			} else {
				// if this is the first error we got, it comes from libvirt and not from a signal:
				// looks like the operation failed, and wasn't aborted.
				if errCode == ErrorCodeNone {
					errCode = ErrorCodeMigrationFailed
				}
			}
		case err = <-monitorError:
			// no implicit abort: it is up to the monitoring code to abort the migration if wishes so.
			if err == nil {
				h.Log().Printf("monitor stop")
			} else {
				h.Log().Printf("monitor stop err=%v", err)
			}
		}
	}
	h.Log().Printf("migration stop errCode=%v", errCode)

	h.Exit(errCode, err)
}

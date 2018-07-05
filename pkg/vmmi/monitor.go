package vmmi

import "io"

type Monitor interface {
	Configure(r io.Reader) error
	Run(resChan chan error)
	Stop()
	Status() (interface{}, error)
}

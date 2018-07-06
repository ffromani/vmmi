package vmmi

import (
	"github.com/fromanirh/vmmi/pkg/vmmi/messages"
	"io"
)

type Monitor interface {
	Configure(r io.Reader) error
	Run(resChan chan error)
	Stop()
	Status(msg *messages.Status) (interface{}, error)
}

package vmmi

type Migrator interface {
	Run(resChan chan error)
}

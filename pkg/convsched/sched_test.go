package convsched

import (
	"os"
	"testing"
)

func TestLoadFromConfiguration(t *testing.T) {
	confFile := "schedconf.json"
	f, err := os.Open(confFile)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}
	defer f.Close()

	cs, err := LoadFromConfiguration(f)

	t.Errorf("sched: %#v", cs)
}

package convsched

import (
	"os"
	"testing"
)

func TestLoadConfiguration(t *testing.T) {
	confFile := "conf_example_simple.json"
	f, err := os.Open(confFile)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}
	defer f.Close()

	conf, err := LoadConfiguration(f)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}

	if len(conf.Schedule.Init) != 1 {
		t.Errorf("sched: %#v", conf.Schedule.Init)
	}
	action := []ConvergenceAction{
		ConvergenceAction{
			Name:   "setDowntime",
			Params: []string{"100"},
		},
	}
	if !convergenceActionSliceEqual(conf.Schedule.Init, action) {
		t.Errorf("sched: %#v", conf.Schedule.Init)
	}
}

func TestHasPostcopy(t *testing.T) {
	checkHasPostcopy(t, "conf_example_simple.json", false)
}

func checkHasPostcopy(t *testing.T, confFile string, expected bool) {
	f, err := os.Open(confFile)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}
	defer f.Close()

	conf, err := LoadConfiguration(f)
	if err != nil {
		t.Errorf("unable to load %v: %v", confFile, err)
	}

	if conf.Schedule.HasPostcopy() != expected {
		if !expected {
			t.Errorf("%v has postcopy, but it should not", confFile)
		} else {
			t.Errorf("%v has not postcopy, but it should", confFile)
		}
	}
}

func convergenceActionSliceEqual(a, b []ConvergenceAction) bool {
	if len(a) != len(b) {
		return false
	}
	for i, val := range a {
		if !convergenceActionEqual(val, b[i]) {
			return false
		}
	}
	return true
}

func convergenceActionEqual(a, b ConvergenceAction) bool {
	if a.Name != b.Name {
		return false
	}
	if len(a.Params) != len(b.Params) {
		return false
	}
	for i, val := range a.Params {
		if val != b.Params[i] {
			return false
		}
	}
	return true
}

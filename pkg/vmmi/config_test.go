package vmmi

import (
	"io/ioutil"
	"testing"
)

func TestFindPluginConfigurationPath(t *testing.T) {
	res := ""
	res = findPluginConfigurationPath([]string{"/some/nested/testing_path"})
	if res != BaseConfigurationDir+"testing_path" {
		t.Errorf("unexpected result: %v", res)
	}
	res = findPluginConfigurationPath([]string{"short_testing_path"})
	if res != BaseConfigurationDir+"short_testing_path" {
		t.Errorf("unexpected result: %v", res)
	}
}

func TestParseParametersFails(t *testing.T) {
	exename := "testingexe"
	exitCode := 0
	term := func(code int) {
		exitCode = code
	}
	h := newHelper()
	h.exitFunc = term
	h.errsink = ioutil.Discard
	h.outsink = ioutil.Discard
	h.parseParameters([]string{exename})

	if exitCode == 0 {
		t.Errorf("readConfiguration unexpected result: %v", exitCode)
	}
}

func TestParseParameters(t *testing.T) {
	exename := "testingexe"
	h := newHelper()
	h.parseParameters([]string{exename, "vmid", "duri", "muri"})
	if h.params.VMid != "vmid" || h.params.URI.Destination != "duri" || h.params.URI.Migration != "muri" || h.params.PluginConfigurationPath != BaseConfigurationDir+exename {
		t.Errorf("unexepcted parameters: %v", h.params)
	}
	h.parseParameters([]string{exename, "vmid", "duri", "muri", "-"})
	if h.params.PluginConfigurationPath != "-" {
		t.Errorf("unexepcted parameters: %v", h.params)
	}
}

func TestReadConfigurationFails(t *testing.T) {
	exitCode := 0
	term := func(code int) {
		exitCode = code
	}
	h := newHelper()
	h.exitFunc = term
	h.errsink = ioutil.Discard
	h.outsink = ioutil.Discard

	h.params.PluginConfigurationPath = "/does/not/exist"
	h.readConfiguration()
	if exitCode == 0 {
		t.Errorf("readConfiguration unexpected result: %v", exitCode)
	}
	if h.confData != nil {
		t.Errorf("readConfiguration unexpected content: %v", h.confData)
	}
}

func TestReadConfiguration(t *testing.T) {
	exitCode := 0
	term := func(code int) {
		exitCode = code
	}
	h := newHelper()
	h.exitFunc = term
	h.errsink = ioutil.Discard
	h.outsink = ioutil.Discard

	h.params.PluginConfigurationPath = "conf_example_simple.json"
	h.readConfiguration()
	if exitCode != 0 {
		t.Errorf("readConfiguration unexpected result: %v", exitCode)
	}
	if h.confData == nil || h.confData.Size() == 0 {
		t.Errorf("readConfiguration unexpected content: %v", h.confData)
	}
}

func TestParseConfiguration(t *testing.T) {
	exitCode := 0
	term := func(code int) {
		exitCode = code
	}
	h := newHelper()
	h.exitFunc = term
	h.errsink = ioutil.Discard
	h.outsink = ioutil.Discard

	h.params.PluginConfigurationPath = "conf_example_simple.json"
	h.readConfiguration()
	h.parseConfiguration()
	if exitCode != 0 {
		t.Errorf("parseConfiguration unexpected result: %v", exitCode)
	}
}

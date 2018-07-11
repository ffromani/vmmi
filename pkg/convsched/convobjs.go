package convsched

import (
	"encoding/json"
	"io"
	"time"
)

const (
	ActionAbort          = "abort"
	ActionEnablePostCopy = "postcopy"
	ActionSetDowntime    = "setDowntime"
)

type ConvergenceAction struct {
	Name   string   `json:"name"`
	Params []string `json:"params"`
}

type ConvergenceItem struct {
	Action ConvergenceAction `json:"action"`
	Limit  int64             `json:"limit"`
}

type ConvergenceSchedule struct {
	Init     []ConvergenceAction `json:"init"`
	Stalling []ConvergenceItem   `json:"stalling"`
}

func (cs *ConvergenceSchedule) HasPostcopy() bool {
	for _, item := range cs.Stalling {
		if item.Action.Name == ActionEnablePostCopy {
			return true
		}
	}
	return false
}

func Load(r io.Reader) (*ConvergenceSchedule, error) {
	dec := json.NewDecoder(r)
	var cs ConvergenceSchedule
	err := dec.Decode(&cs)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

type ConvergenceScheduleConfiguration struct {
	Schedule        ConvergenceSchedule `json:"schedule"`
	MonitorInterval time.Duration       `json:"monitorInterval"`
}

type ConfigurationMessage struct {
	Configuration ConvergenceScheduleConfiguration `json:"configuration"`
}

func LoadConfiguration(r io.Reader) (*ConvergenceScheduleConfiguration, error) {
	dec := json.NewDecoder(r)
	var conf ConfigurationMessage
	err := dec.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf.Configuration, nil
}

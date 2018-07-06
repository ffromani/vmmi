package convsched

import (
	"encoding/json"
	"io"
)

type ConvergenceAction struct {
	Name   string   `json:"name"`
	Params []string `json:"params"`
}

type ConvergenceItem struct {
	Action ConvergenceAction `json:"action"`
	Limit  int               `json:"limit"`
}

type ConvergenceSchedule struct {
	Init     []ConvergenceAction `json:"init"`
	Stalling []ConvergenceItem   `json:"stalling"`
}

func (cs *ConvergenceSchedule) HasPostcopy() bool {
	return true
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
	Schedule ConvergenceSchedule `json:"schedule"`
}

type ConfigurationMessage struct {
	Configuration ConvergenceScheduleConfiguration `json:"configuration"`
}

func LoadFromConfiguration(r io.Reader) (*ConvergenceSchedule, error) {
	dec := json.NewDecoder(r)
	var conf ConfigurationMessage
	err := dec.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return &conf.Configuration.Schedule, nil
}

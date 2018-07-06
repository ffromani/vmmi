package messages

import (
	"bytes"
	"encoding/json"
	"log"
	"testing"
)

func TestReportSuccess(t *testing.T) {
	var W bytes.Buffer
	var L bytes.Buffer
	sink := Sink{
		W: &W,
		L: log.New(&L, "test: ", log.LstdFlags),
	}
	sink.ReportSuccess()

	var msg CompletionSuccess
	dec := json.NewDecoder(&W)
	dec.Decode(&msg)

	if msg.Version != Version {
		t.Errorf("wrong version: %v vs %v", msg.Version, Version)
	}
	if msg.ContentType != ContentTypeCompletion {
		t.Errorf("wrong contentType: %v vs %v", msg.ContentType, ContentTypeCompletion)
	}
	if msg.Timestamp <= 0 {
		t.Errorf("wrong timestamp: %v", msg.Timestamp)
	}

	if msg.Completion.Result != CompletionResultSuccess {
		t.Errorf("wrong completion result: %v vs %v", msg.Completion.Result, CompletionResultSuccess)
	}
}

func TestReportError(t *testing.T) {
	var W bytes.Buffer
	var L bytes.Buffer
	sink := Sink{
		W: &W,
		L: log.New(&L, "test: ", log.LstdFlags),
	}
	errData := ErrorData{
		Code:    4242,
		Message: "testing error message",
		Details: "testing error message details, slightly more verbose",
	}
	sink.ReportError(errData.Code, errData.Message, errData.Details)

	var msg CompletionError
	dec := json.NewDecoder(&W)
	dec.Decode(&msg)

	if msg.Version != Version {
		t.Errorf("wrong version: %v vs %v", msg.Version, Version)
	}
	if msg.ContentType != ContentTypeCompletion {
		t.Errorf("wrong contentType: %v vs %v", msg.ContentType, ContentTypeCompletion)
	}
	if msg.Timestamp <= 0 {
		t.Errorf("wrong timestamp: %v", msg.Timestamp)
	}

	if msg.Completion.Result != CompletionResultError {
		t.Errorf("wrong completion result: %v vs %v", msg.Completion.Result, CompletionResultError)
	}
	if *msg.Completion.Error != errData {
		t.Errorf("wrong completion data: %v vs %v", msg.Completion.Error, errData)
	}
}

type statusData struct {
	ProgressPercentage uint64 `json:"percentage"`
}

type testStatusMessage struct {
	Status
	StatusData statusData `json:"status"`
}

func TestReportStatusPassthrough(t *testing.T) {
	var W bytes.Buffer
	var L bytes.Buffer
	sink := Sink{
		W: &W,
		L: log.New(&L, "test: ", log.LstdFlags),
	}
	sink.ReportStatus(NewStatus())

	var msg Status
	dec := json.NewDecoder(&W)
	dec.Decode(&msg)

	if msg.Version != Version {
		t.Errorf("wrong version: %v vs %v", msg.Version, Version)
	}
	if msg.ContentType != ContentTypeStatus {
		t.Errorf("wrong contentType: %v vs %v", msg.ContentType, ContentTypeStatus)
	}
	if msg.Timestamp <= 0 {
		t.Errorf("wrong timestamp: %v", msg.Timestamp)
	}
}

func TestReportStatusAugmented(t *testing.T) {
	var W bytes.Buffer
	var L bytes.Buffer
	sink := Sink{
		W: &W,
		L: log.New(&L, "test: ", log.LstdFlags),
	}

	payload := statusData{
		ProgressPercentage: 51,
	}
	st := testStatusMessage{
		Status:     *NewStatus(),
		StatusData: payload,
	}
	sink.ReportStatus(st)

	var msg testStatusMessage
	dec := json.NewDecoder(&W)
	dec.Decode(&msg)

	if msg.Version != Version {
		t.Errorf("wrong version: %v vs %v", msg.Version, Version)
	}
	if msg.ContentType != ContentTypeStatus {
		t.Errorf("wrong contentType: %v vs %v", msg.ContentType, ContentTypeStatus)
	}
	if msg.Timestamp <= 0 {
		t.Errorf("wrong timestamp: %v", msg.Timestamp)
	}
	if msg.StatusData != payload {
		t.Errorf("wrong payload: %v vs %v", msg.StatusData, payload)
	}
}

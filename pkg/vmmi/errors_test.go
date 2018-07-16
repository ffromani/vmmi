package vmmi

import "testing"

func TestStrerror(t *testing.T) {
	msg := ""

	msg = Strerror(-1)
	if msg != "unknown" {
		t.Errorf("unexpected error message %v for code %v", msg, -1)
	}
	msg = Strerror(0)
	if msg != "" {
		t.Errorf("unexpected error message %v for code %v", msg, 0)
	}

	for i := 1; i < ErrorCodeLast; i++ {
		msg = Strerror(i)
		if msg == "" || msg == "unknown" {
			t.Errorf("unexpected error message %v for code %v", msg, i)
		}
	}
}

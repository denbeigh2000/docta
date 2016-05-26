package checkers

import (
	"testing"

	"github.com/denbeigh2000/docta"
)

func TestCheckStringGreen(t *testing.T) {
	okText := "This is OK"

	state := CheckString(okText, "not ok", "definitely not ok")
	if state != docta.Green {
		t.Errorf("Got state %v, expected %v", state, docta.Green)
	}
}

func TestCheckStringYellow(t *testing.T) {
	yellowText := "This is not OK"

	state := CheckString(yellowText, "This is not OK", "definitely not ok")
	if state != docta.Yellow {
		t.Errorf("Got state %v, expected %v", state, docta.Yellow)
	}
}

func TestCheckStringRed(t *testing.T) {
	yellowText := "This is definitely not OK"

	state := CheckString(yellowText, "not ok", "This is definitely not OK")
	if state != docta.Red {
		t.Errorf("Got state %v, expected %v", state, docta.Red)
	}
}

package docta

import (
	"testing"

	"strings"
)

type fakeChecker State

const (
	greenMessage  = "Shit is OK"
	yellowMessage = "Shit might be borked"
	redMessage    = "Shit be borked"
)

func (c fakeChecker) Check() HealthState {
	var info string

	switch State(c) {
	case Green:
		info = greenMessage
	case Yellow:
		info = yellowMessage
	case Red:
		info = redMessage
	default:
		panic("Invalid state")
	}

	return HealthState{State(c), info}
}

func TestMultiCheckBasic(t *testing.T) {
	checkers := []Checker{
		fakeChecker(Green), fakeChecker(Green),
		fakeChecker(Red), fakeChecker(Red),
		fakeChecker(Red), fakeChecker(Yellow),
	}

	multiChecker := MultiChecker{checkers}

	state := multiChecker.Check()
	splitStates := strings.Split(state.Info, SplitCharacter)
	if len(splitStates) != len(checkers) {
		t.Errorf("Expected %v messages, got %v")
	}

	infoMap := make(map[string]int)
	for _, s := range splitStates {
		rec, ok := infoMap[s]
		if !ok {
			infoMap[s] = 1
		} else {
			infoMap[s] = rec + 1
		}
	}

	if infoMap[greenMessage] != 2 {
		t.Errorf("Expected %v occurances of %v, got %v", 2, greenMessage, infoMap[greenMessage])
	}

	if infoMap[yellowMessage] != 1 {
		t.Errorf("Expected %v occurances of %v, got %v", 2, yellowMessage, infoMap[yellowMessage])
	}

	if infoMap[redMessage] != 3 {
		t.Errorf("Expected %v occurances of %v, got %v", 3, redMessage, infoMap[redMessage])
	}
}

func TestMultiCheckGreen(t *testing.T) {
	checkers := []Checker{fakeChecker(Green), fakeChecker(Green)}
	multiChecker := MultiChecker{checkers}

	state := multiChecker.Check()
	if state.State != Green {
		t.Errorf("Expected state to be %v, got %v", Green, state.State)
	}
}

func TestMultiCheckYellow(t *testing.T) {
	multiChecker := MultiChecker{[]Checker{fakeChecker(Green), fakeChecker(Yellow)}}
	state := multiChecker.Check()
	if state.State != Yellow {
		t.Errorf("Expected state to be %v, got %v", Yellow, state.State)
	}
}

func TestMultiCheckRed(t *testing.T) {
	multiChecker := MultiChecker{[]Checker{fakeChecker(Red), fakeChecker(Yellow)}}

	state := multiChecker.Check()

	if state.State != Red {
		t.Errorf("Expected state to be %v, got %v", Red, state.State)
	}
}

func TestNoChecks(t *testing.T) {
	multiChecker := MultiChecker{}
	state := multiChecker.Check()
	if state.Info != DefaultInfo {
		t.Errorf("Expected info to be %v, got %v", DefaultInfo, state.Info)
	}

	if state.State != DefaultState {
		t.Errorf("Expected state to be %v, got %v", DefaultState, state.State)
	}
}

func TestService(t *testing.T) {
	checkers := map[string]Checker{
		"red":   fakeChecker(Red),
		"green": fakeChecker(Green),
	}
	service := NewService("Sample Service", checkers)
	service.Check()

	if _, ok := service.CurrentStatus["red"]; !ok {
		t.Errorf("Missing state \"%v\"", "red")
	}
	if _, ok := service.CurrentStatus["green"]; !ok {
		t.Errorf("Missing state \"%v\"", "green")
	}

	expectedGreenState := HealthState{Green, greenMessage}
	greenState := service.CurrentStatus["green"]
	if greenState != expectedGreenState {
		t.Errorf("Green state not recorded properly - got %v", greenState)
	}

	expectedRedState := HealthState{Red, redMessage}
	redState := service.CurrentStatus["red"]
	if redState != expectedRedState {
		t.Errorf("Red state not recorded properly - got %v", redState)
	}
}

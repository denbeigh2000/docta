package docta

import (
	"strings"
	"sync"
)

type State int

const (
	Green State = iota
	Yellow
	Red
)

const (
	DefaultState   = Green
	DefaultInfo    = "OK"
	SplitCharacter = ";"
)

type HealthState struct {
	State State
	Info  string
}

type Checker interface {
	Check() HealthState
}

type MultiChecker struct {
	Checkers []Checker
}

func (c MultiChecker) Check() HealthState {
	state := DefaultState
	checks := make(chan HealthState)

	var infoMessages []string

	checkGroup := sync.WaitGroup{}
	checkGroup.Add(len(c.Checkers))
	go func() {
		checkGroup.Wait()
		close(checks)
	}()

	for _, checker := range c.Checkers {
		go func(checker Checker) {
			checks <- checker.Check()
			checkGroup.Done()
		}(checker)
	}

	for result := range checks {
		if result.Info != "" {
			infoMessages = append(infoMessages, result.Info)
		}

		if result.State > state {
			state = result.State
		}
	}

	var info string
	if len(infoMessages) > 0 {
		info = strings.Join(infoMessages, SplitCharacter)
	} else {
		info = DefaultInfo
	}

	return HealthState{state, info}
}

func NewService(friendlyName string, checkers map[string]Checker) Service {
	return Service{
		FriendlyName:  friendlyName,
		CurrentStatus: make(map[string]HealthState),
		Checkers:      checkers,
	}
}

// A Service represents a top-level service that has no overall state, but
// has multiple sub-states
type Service struct {
	FriendlyName  string                 `json:"friendly_name"`
	CurrentStatus map[string]HealthState `json:"current_status"`
	Checkers      map[string]Checker     `json:"-"`
}

func (s *Service) Check() {
	wg := sync.WaitGroup{}
	wg.Add(len(s.Checkers))

	for key, checker := range s.Checkers {
		go func(key string, checker Checker) {
			s.CurrentStatus[key] = checker.Check()
			wg.Done()
		}(key, checker)
	}

	wg.Wait()
}

type Docta struct {
	Services map[string]Service `json:"services"`
}

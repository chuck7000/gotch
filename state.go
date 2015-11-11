package gotch

import (
	"sync"
)

// State provides the major state tracking fields and enforces mutex protections on them.
type State struct {
	sync.RWMutex
	quietPeriod    bool
	running        bool
	changeHappened bool
}

// Quiet gets the quietPeriod value
func (s *State) Quiet() bool {
	s.RLock()
	defer s.RUnlock()
	return s.quietPeriod
}

// SetQuiet sets the quietPeriod value
func (s *State) SetQuiet(v bool) {
	s.Lock()
	s.quietPeriod = v
	s.Unlock()
}

// Running gets the running value
func (s *State) Running() bool {
	s.RLock()
	defer s.RUnlock()
	return s.running
}

// SetRunning sets the running value
func (s *State) SetRunning(v bool) {
	s.Lock()
	s.running = v
	s.Unlock()
}

// ChangeHappened sets the changeHappened value
func (s *State) ChangeHappened() bool {
	s.RLock()
	defer s.RUnlock()
	return s.changeHappened
}

// SetChangeHappened gets the changeHappened value
func (s *State) SetChangeHappened(v bool) {
	s.Lock()
	s.changeHappened = v
	s.Unlock()
}

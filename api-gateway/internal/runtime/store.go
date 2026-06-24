package runtime

import "sync"

type StateStore struct {
	mu sync.RWMutex

	states map[string]State
}

func NewStateStore() *StateStore {
	return &StateStore{
		states: make(map[string]State),
	}
}

func (s *StateStore) Set(name string, state State) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.states[name] = state
}

func (s *StateStore) Get(component string) State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state, ok := s.states[component]
	if !ok {
		return StateStopped
	}

	return state
}

package agent

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu     sync.RWMutex
	agents map[string]Agent
}

func NewRegistry() *Registry {
	return &Registry{
		agents: make(map[string]Agent),
	}
}

func (r *Registry) Register(a Agent) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	name := a.Name()
	if _, exists := r.agents[name]; exists {
		return fmt.Errorf("agent %q already registered", name)
	}
	r.agents[name] = a
	return nil
}

func (r *Registry) Get(name string) (Agent, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.agents[name]
	if !ok {
		return nil, fmt.Errorf("agent %q not found", name)
	}
	return a, nil
}

func (r *Registry) List() []Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]Agent, 0, len(r.agents))
	for _, a := range r.agents {
		result = append(result, a)
	}
	return result
}

func (r *Registry) Available() []Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Agent
	for _, a := range r.agents {
		if a.IsAvailable() {
			result = append(result, a)
		}
	}
	return result
}

func (r *Registry) FindBySpecialty(task TaskType) []Agent {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []Agent
	for _, a := range r.agents {
		for _, s := range a.Specialties() {
			if s == task {
				result = append(result, a)
				break
			}
		}
	}
	return result
}

type AgentFactory func() (Agent, error)

var (
	factoryMu  sync.RWMutex
	factories  = make(map[string]AgentFactory)
)

func RegisterFactory(name string, f AgentFactory) {
	factoryMu.Lock()
	defer factoryMu.Unlock()
	factories[name] = f
}

func CreateAgent(name string) (Agent, error) {
	factoryMu.RLock()
	f, ok := factories[name]
	factoryMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("no factory registered for agent %q", name)
	}
	return f()
}

func RegisteredFactories() []string {
	factoryMu.RLock()
	defer factoryMu.RUnlock()
	names := make([]string, 0, len(factories))
	for name := range factories {
		names = append(names, name)
	}
	return names
}

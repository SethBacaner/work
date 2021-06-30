package internal

import (
	"errors"
	"sync"
)

type Registry struct {
	definitions map[string]*TaskDefinition
	// TODO: can be a read write mutex
	mu *sync.Mutex
}

func (r *Registry) Register(definition *TaskDefinition) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.definitions[definition.Name] = definition
}

func (r *Registry) DefinitionFor(name string) (*TaskDefinition, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	def, ok := r.definitions[name]
	if !ok {
		return nil, errors.New("no definition found")
	}

	return def, nil
}

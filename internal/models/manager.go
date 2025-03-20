package models

import (
	"fmt"
	"sync"
)

// Manager handles model management operations
type Manager struct {
	mu sync.RWMutex
	// models maps model IDs to their respective models
	models map[string]*ONNXModel
}

// NewManager creates a new model manager
func NewManager() *Manager {
	return &Manager{
		models: make(map[string]*ONNXModel),
	}
}

// GetModel retrieves a model by its ID
func (m *Manager) GetModel(id string) (*ONNXModel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	model, exists := m.models[id]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", id)
	}

	return model, nil
}

// ListModels returns all available models
func (m *Manager) ListModels() ([]*ONNXModel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	models := make([]*ONNXModel, 0, len(m.models))
	for _, model := range m.models {
		models = append(models, model)
	}

	return models, nil
}

// AddModel adds a new model to the manager
func (m *Manager) AddModel(model *ONNXModel) error {
	if model == nil {
		return fmt.Errorf("model cannot be nil")
	}
	if model.ID == "" {
		return fmt.Errorf("model ID cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.models[model.ID]; exists {
		return fmt.Errorf("model already exists: %s", model.ID)
	}

	m.models[model.ID] = model
	return nil
}

// RemoveModel removes a model from the manager
func (m *Manager) RemoveModel(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.models[id]; !exists {
		return fmt.Errorf("model not found: %s", id)
	}

	delete(m.models, id)
	return nil
}

// UpdateModel updates an existing model
func (m *Manager) UpdateModel(model *ONNXModel) error {
	if model == nil {
		return fmt.Errorf("model cannot be nil")
	}
	if model.ID == "" {
		return fmt.Errorf("model ID cannot be empty")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.models[model.ID]; !exists {
		return fmt.Errorf("model not found: %s", model.ID)
	}

	m.models[model.ID] = model
	return nil
}

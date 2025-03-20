package inference

import (
	"context"
	"fmt"
	"sync"

	"github.com/josealecrim/audiototext/internal/hardware"
	"github.com/josealecrim/audiototext/internal/models"
	"github.com/pkg/errors"
)

// Manager handles ONNX Runtime sessions and their lifecycle
type Manager struct {
	// mu protects the sessions map
	mu sync.RWMutex
	// sessions maps model IDs to their respective sessions
	sessions map[string]*Session
	// hwDetector is used to determine available hardware capabilities
	hwDetector *hardware.Detector
	// defaultConfig is the default session configuration
	defaultConfig SessionConfig
	// defaultBatchConfig is the default batch processing configuration
	defaultBatchConfig BatchConfig
}

// NewManager creates a new session manager
func NewManager(hwDetector *hardware.Detector) *Manager {
	defaultConfig := SessionConfig{
		ExecutionProvider:      CPUExecutionProvider,
		InterOpNumThreads:      4,
		IntraOpNumThreads:      4,
		GraphOptimizationLevel: 99, // Maximum optimization
		EnableMemoryPattern:    true,
		EnableCPUMemArena:      true,
		EnableProfiling:        false,
	}

	defaultBatchConfig := BatchConfig{
		MaxBatchSize:        32,
		DynamicBatching:     true,
		MaxLatencyMs:        100,
		PreferredBatchSizes: []int{1, 2, 4, 8, 16, 32},
	}

	// Adjust configuration based on hardware capabilities
	if hwDetector.HasCUDAGPU() {
		defaultConfig.ExecutionProvider = CUDAExecutionProvider
	} else if hwDetector.HasIntelGPU() {
		defaultConfig.ExecutionProvider = OpenVINOExecutionProvider
	}

	// Adjust thread counts based on CPU cores
	cpuInfo := hwDetector.GetCPUInfo()
	if cpuInfo.Cores > 2 {
		defaultConfig.InterOpNumThreads = cpuInfo.Cores / 2
		defaultConfig.IntraOpNumThreads = cpuInfo.Cores / 2
	}

	return &Manager{
		sessions:           make(map[string]*Session),
		hwDetector:         hwDetector,
		defaultConfig:      defaultConfig,
		defaultBatchConfig: defaultBatchConfig,
	}
}

// CreateSession creates a new ONNX Runtime session for the given model
func (m *Manager) CreateSession(ctx context.Context, model *models.ONNXModel, config *SessionConfig, batchConfig *BatchConfig) (*Session, error) {
	if model == nil {
		return nil, errors.New("model cannot be nil")
	}

	// Use default configurations if not provided
	if config == nil {
		config = &m.defaultConfig
	}
	if batchConfig == nil {
		batchConfig = &m.defaultBatchConfig
	}

	// Create session
	session := &Session{
		Model:       model,
		Config:      *config,
		BatchConfig: *batchConfig,
	}

	// Initialize ONNX Runtime session
	if err := m.initializeSession(ctx, session); err != nil {
		return nil, errors.Wrap(err, "failed to initialize session")
	}

	// Store session
	m.mu.Lock()
	m.sessions[model.ID] = session
	m.mu.Unlock()

	return session, nil
}

// GetSession retrieves an existing session for the given model ID
func (m *Manager) GetSession(modelID string) (*Session, error) {
	m.mu.RLock()
	session, exists := m.sessions[modelID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no session found for model ID: %s", modelID)
	}

	return session, nil
}

// CloseSession closes and removes a session for the given model ID
func (m *Manager) CloseSession(modelID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.sessions[modelID]
	if !exists {
		return fmt.Errorf("no session found for model ID: %s", modelID)
	}

	// Close the session
	if err := m.closeSession(session); err != nil {
		return errors.Wrap(err, "failed to close session")
	}

	// Remove from map
	delete(m.sessions, modelID)
	return nil
}

// CloseAllSessions closes all active sessions
func (m *Manager) CloseAllSessions() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var errs []error
	for modelID, session := range m.sessions {
		if err := m.closeSession(session); err != nil {
			errs = append(errs, fmt.Errorf("failed to close session for model %s: %v", modelID, err))
		}
		delete(m.sessions, modelID)
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing sessions: %v", errs)
	}
	return nil
}

// GetStats returns statistics for all active sessions
func (m *Manager) GetStats() map[string]Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()

	stats := make(map[string]Stats)
	for modelID, session := range m.sessions {
		// Get stats from session
		sessionStats := m.getSessionStats(session)
		stats[modelID] = sessionStats
	}

	return stats
}

// initializeSession initializes the ONNX Runtime session
func (m *Manager) initializeSession(ctx context.Context, session *Session) error {
	// TODO: Implement actual ONNX Runtime session initialization
	// This will be implemented when we integrate the ONNX Runtime Go bindings
	return nil
}

// closeSession closes an ONNX Runtime session
func (m *Manager) closeSession(session *Session) error {
	// TODO: Implement actual ONNX Runtime session cleanup
	// This will be implemented when we integrate the ONNX Runtime Go bindings
	return nil
}

// getSessionStats retrieves statistics for a session
func (m *Manager) getSessionStats(session *Session) Stats {
	// TODO: Implement actual statistics collection from ONNX Runtime session
	return Stats{}
}

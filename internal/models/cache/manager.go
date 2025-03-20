package cache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/josealecrim/audiototext/internal/models"
)

// Manager implementa o gerenciador de cache
type Manager struct {
	config models.Config
	mu     sync.RWMutex
	index  map[string]models.ModelInfo
}

// NewManager cria um novo gerenciador de cache
func NewManager(config models.Config) (*Manager, error) {
	m := &Manager{
		config: config,
		index:  make(map[string]models.ModelInfo),
	}

	if err := m.loadIndex(); err != nil {
		return nil, err
	}

	return m, nil
}

// Store armazena um modelo no cache
func (m *Manager) Store(model *models.Model) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Atualiza o índice
	m.index[model.Info.ID] = model.Info

	// Salva o índice
	return m.saveIndex()
}

// Load carrega um modelo do cache
func (m *Manager) Load(modelID string) (*models.Model, error) {
	m.mu.RLock()
	info, exists := m.index[modelID]
	m.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("modelo não encontrado no cache: %s", modelID)
	}

	path := filepath.Join(m.config.CachePath, modelID)
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("arquivo do modelo não encontrado: %w", err)
	}

	return &models.Model{
		Info: info,
		Path: path,
	}, nil
}

// List lista todos os modelos no cache
func (m *Manager) List() []models.ModelInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var infos []models.ModelInfo
	for _, info := range m.index {
		infos = append(infos, info)
	}

	// Ordena por data de download
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Downloaded.After(infos[j].Downloaded)
	})

	return infos
}

// Clean limpa o cache de acordo com a política
func (m *Manager) Clean(policy models.CleanPolicy) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var toDelete []string
	var totalSize int64

	// Lista todos os modelos ordenados por data de último uso
	var infos []models.ModelInfo
	for _, info := range m.index {
		infos = append(infos, info)
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].LastUsed.Before(infos[j].LastUsed)
	})

	// Mantém os N modelos mais recentes
	if len(infos) > policy.KeepLatest {
		infos = infos[policy.KeepLatest:]
	}

	now := time.Now()
	for _, info := range infos {
		// Verifica idade do modelo
		if policy.MaxAge > 0 && now.Sub(info.LastUsed) > policy.MaxAge {
			toDelete = append(toDelete, info.ID)
			continue
		}

		// Verifica tamanho total
		totalSize += info.Size
		if policy.MaxSize > 0 && totalSize > policy.MaxSize {
			toDelete = append(toDelete, info.ID)
		}
	}

	// Remove os modelos marcados para deleção
	for _, id := range toDelete {
		if err := m.remove(id); err != nil {
			return fmt.Errorf("erro ao remover modelo %s: %w", id, err)
		}
	}

	return m.saveIndex()
}

// remove remove um modelo do cache
func (m *Manager) remove(modelID string) error {
	path := filepath.Join(m.config.CachePath, modelID)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	delete(m.index, modelID)
	return nil
}

// loadIndex carrega o índice do cache
func (m *Manager) loadIndex() error {
	path := filepath.Join(m.config.CachePath, "index.json")
	data, err := ioutil.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("erro ao ler índice: %w", err)
	}

	return json.Unmarshal(data, &m.index)
}

// saveIndex salva o índice do cache
func (m *Manager) saveIndex() error {
	data, err := json.MarshalIndent(m.index, "", "    ")
	if err != nil {
		return fmt.Errorf("erro ao serializar índice: %w", err)
	}

	path := filepath.Join(m.config.CachePath, "index.json")
	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("erro ao salvar índice: %w", err)
	}

	return nil
}

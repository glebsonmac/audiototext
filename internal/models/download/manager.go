package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/josealecrim/audiototext/internal/models"
)

// Manager implementa o gerenciador de downloads
type Manager struct {
	config    models.Config
	client    *http.Client
	downloads map[string]*download
	mu        sync.RWMutex
}

type download struct {
	progress models.Progress
	cancel   context.CancelFunc
}

// NewManager cria um novo gerenciador de downloads
func NewManager(config models.Config) *Manager {
	return &Manager{
		config:    config,
		client:    &http.Client{Timeout: config.DownloadTimeout},
		downloads: make(map[string]*download),
	}
}

// Download inicia o download de um modelo
func (m *Manager) Download(modelType models.ModelType, version string) (*models.Model, error) {
	modelID := fmt.Sprintf("%s-%s", modelType, version)

	// Verifica se já existe um download em andamento
	m.mu.RLock()
	if _, exists := m.downloads[modelID]; exists {
		m.mu.RUnlock()
		return nil, fmt.Errorf("download já em andamento para %s", modelID)
	}
	m.mu.RUnlock()

	// Cria diretório de cache se não existir
	if err := os.MkdirAll(m.config.CachePath, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de cache: %w", err)
	}

	// Prepara o download
	ctx, cancel := context.WithCancel(context.Background())
	dl := &download{
		progress: models.Progress{Status: "iniciando"},
		cancel:   cancel,
	}

	m.mu.Lock()
	m.downloads[modelID] = dl
	m.mu.Unlock()

	// Inicia o download em uma goroutine
	go func() {
		defer func() {
			m.mu.Lock()
			delete(m.downloads, modelID)
			m.mu.Unlock()
			cancel()
		}()

		url := getModelURL(modelType, version)
		if err := m.downloadFile(ctx, url, modelID); err != nil {
			dl.progress.Status = fmt.Sprintf("erro: %v", err)
			return
		}

		dl.progress.Status = "concluído"
	}()

	// Retorna imediatamente com as informações do modelo
	return &models.Model{
		Info: models.ModelInfo{
			ID:         modelID,
			Type:       modelType,
			Version:    version,
			Downloaded: time.Now(),
		},
		Path: filepath.Join(m.config.CachePath, modelID),
	}, nil
}

// GetProgress retorna o progresso de um download
func (m *Manager) GetProgress(modelID string) (models.Progress, error) {
	m.mu.RLock()
	dl, exists := m.downloads[modelID]
	m.mu.RUnlock()

	if !exists {
		return models.Progress{}, fmt.Errorf("download não encontrado: %s", modelID)
	}

	return dl.progress, nil
}

// Cancel cancela um download em andamento
func (m *Manager) Cancel(modelID string) error {
	m.mu.Lock()
	dl, exists := m.downloads[modelID]
	m.mu.Unlock()

	if !exists {
		return fmt.Errorf("download não encontrado: %s", modelID)
	}

	dl.cancel()
	return nil
}

// downloadFile realiza o download de um arquivo
func (m *Manager) downloadFile(ctx context.Context, url, modelID string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}

	resp, err := m.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao fazer request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code inválido: %d", resp.StatusCode)
	}

	file, err := os.Create(filepath.Join(m.config.CachePath, modelID))
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo: %w", err)
	}
	defer file.Close()

	// Atualiza o progresso durante o download
	m.mu.Lock()
	dl := m.downloads[modelID]
	dl.progress.Total = resp.ContentLength
	m.mu.Unlock()

	buffer := make([]byte, 32*1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := resp.Body.Read(buffer)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return fmt.Errorf("erro ao ler resposta: %w", err)
			}

			if _, err := file.Write(buffer[:n]); err != nil {
				return fmt.Errorf("erro ao escrever arquivo: %w", err)
			}

			m.mu.Lock()
			dl.progress.Current += int64(n)
			dl.progress.Percentage = float64(dl.progress.Current) / float64(dl.progress.Total) * 100
			m.mu.Unlock()
		}
	}
}

// getModelURL retorna a URL do modelo baseado no tipo e versão
func getModelURL(modelType models.ModelType, version string) string {
	baseURL := "https://huggingface.co/models"

	switch modelType {
	case models.TypeClaude:
		return fmt.Sprintf("%s/anthropic/claude-3-sonnet-%s/resolve/main/model.onnx", baseURL, version)
	case models.TypeWhisper:
		return fmt.Sprintf("%s/openai/whisper-small-%s/resolve/main/model.onnx", baseURL, version)
	default:
		return ""
	}
}

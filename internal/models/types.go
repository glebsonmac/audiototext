package models

import (
	"time"
)

// ModelType representa o tipo do modelo
type ModelType string

const (
	TypeClaude  ModelType = "claude"
	TypeWhisper ModelType = "whisper"
)

// ModelInfo contém informações sobre um modelo
type ModelInfo struct {
	ID          string    `json:"id"`
	Type        ModelType `json:"type"`
	Version     string    `json:"version"`
	Size        int64     `json:"size"`
	Downloaded  time.Time `json:"downloaded"`
	LastUsed    time.Time `json:"last_used"`
	IsConverted bool      `json:"is_converted"`
}

// Model representa um modelo baixado
type Model struct {
	Info     ModelInfo `json:"info"`
	Path     string    `json:"path"`
	Checksum string    `json:"checksum"`
}

// ONNXModel extends Model to include ONNX-specific fields
type ONNXModel struct {
	*Model
	// OptimizationLevel represents the level of optimization applied to the model
	OptimizationLevel int
	// IsQuantized indicates if the model has been quantized
	IsQuantized bool
	// ID is a unique identifier for the model
	ID string
}

// Progress representa o progresso de uma operação
type Progress struct {
	Total      int64   `json:"total"`
	Current    int64   `json:"current"`
	Percentage float64 `json:"percentage"`
	Status     string  `json:"status"`
}

// CleanPolicy define a política de limpeza de cache
type CleanPolicy struct {
	MaxSize    int64         `json:"max_size"`
	MaxAge     time.Duration `json:"max_age"`
	KeepLatest int           `json:"keep_latest"`
}

// Config contém a configuração do sistema de modelos
type Config struct {
	CachePath           string        `yaml:"cache_path"`
	MaxCacheSize        string        `yaml:"max_cache_size"`
	RetentionPeriod     string        `yaml:"retention_period"`
	ConcurrentDownloads int           `yaml:"concurrent_downloads"`
	DownloadTimeout     time.Duration `yaml:"download_timeout"`
	Retries             int           `yaml:"retries"`
	OptimizationLevel   string        `yaml:"optimization_level"`
	EnableQuantization  bool          `yaml:"enable_quantization"`
	TargetPlatform      string        `yaml:"target_platform"`
}

package models

import "context"

// Manager defines the interface for model management operations
type Manager interface {
	// Model Management
	GetModel(id string) error
	ListModels() ([]string, error)
	AddModel(id string) error
	RemoveModel(id string) error
	UpdateModel(id string) error

	// Model Operations
	DownloadModel(ctx context.Context, modelID string) error
	ModelExists(modelID string) bool
	ConvertToONNX(ctx context.Context, modelID string) (string, error)
	ONNXModelExists(modelID string) bool
	ValidateONNXModel(ctx context.Context, modelID string) error
	OptimizeONNXModel(ctx context.Context, modelID string) (string, error)
	OptimizedModelExists(modelID string) bool
	GetModelVersion(modelID string) (string, error)

	// Backup Operations
	CreateBackup(ctx context.Context, modelID string) error
	ListBackups(modelID string) ([]string, error)
	RestoreFromBackup(ctx context.Context, modelID string, backupID string) error

	// Cache Operations
	GetCacheSize() (int64, error)
	ClearCache() error
}

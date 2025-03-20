package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/josealecrim/audiototext/internal/models"
	"github.com/josealecrim/audiototext/internal/models/cache"
	"github.com/josealecrim/audiototext/internal/models/download"
)

func main() {
	// Configuração
	config := models.Config{
		CachePath:           filepath.Join(".", "models", "cache"),
		MaxCacheSize:        "10GB",
		RetentionPeriod:     "30d",
		ConcurrentDownloads: 2,
		DownloadTimeout:     time.Hour,
		Retries:             3,
		OptimizationLevel:   "O2",
		EnableQuantization:  true,
		TargetPlatform:      "cpu",
	}

	// Cria diretório de cache
	if err := os.MkdirAll(config.CachePath, 0755); err != nil {
		log.Fatalf("Erro ao criar diretório de cache: %v", err)
	}

	// Inicializa gerenciadores
	downloader := download.NewManager(config)
	cacheManager, err := cache.NewManager(config)
	if err != nil {
		log.Fatalf("Erro ao criar gerenciador de cache: %v", err)
	}

	// Lista de modelos para download
	modelsToDownload := []struct {
		Type    models.ModelType
		Version string
	}{
		{models.TypeClaude, "v1"},
		{models.TypeWhisper, "v1"},
	}

	// Download dos modelos
	for _, m := range modelsToDownload {
		fmt.Printf("Baixando modelo %s-%s...\n", m.Type, m.Version)

		// Inicia download
		model, err := downloader.Download(m.Type, m.Version)
		if err != nil {
			log.Printf("Erro ao iniciar download de %s-%s: %v\n", m.Type, m.Version, err)
			continue
		}

		// Monitora progresso
		for {
			progress, err := downloader.GetProgress(model.Info.ID)
			if err != nil {
				log.Printf("Erro ao obter progresso: %v\n", err)
				break
			}

			fmt.Printf("\rProgresso: %.2f%% (%s)", progress.Percentage, progress.Status)

			if progress.Status == "concluído" || progress.Status == "erro" {
				fmt.Println()
				break
			}

			time.Sleep(time.Second)
		}

		// Armazena no cache
		if err := cacheManager.Store(model); err != nil {
			log.Printf("Erro ao armazenar modelo no cache: %v\n", err)
			continue
		}
	}

	// Lista modelos no cache
	fmt.Println("\nModelos no cache:")
	for _, info := range cacheManager.List() {
		fmt.Printf("- %s (v%s): baixado em %s\n",
			info.Type, info.Version, info.Downloaded.Format(time.RFC3339))
	}

	// Limpa cache antigo
	policy := models.CleanPolicy{
		MaxSize:    1024 * 1024 * 1024 * 10, // 10GB
		MaxAge:     30 * 24 * time.Hour,     // 30 dias
		KeepLatest: 5,
	}

	if err := cacheManager.Clean(policy); err != nil {
		log.Printf("Erro ao limpar cache: %v\n", err)
	}
}

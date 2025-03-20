package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/josealecrim/audiototext/internal/hardware/detection"
	"github.com/josealecrim/audiototext/internal/hardware/monitoring"
)

func main() {
	// Cria detector de hardware
	detector := detection.NewDetector()

	// Detecta hardware
	info, err := detector.Detect()
	if err != nil {
		log.Fatalf("Erro ao detectar hardware: %v", err)
	}

	// Imprime informações do hardware em JSON
	jsonInfo, err := json.MarshalIndent(info, "", "    ")
	if err != nil {
		log.Fatalf("Erro ao serializar informações: %v", err)
	}
	fmt.Printf("Informações do Hardware:\n%s\n\n", string(jsonInfo))

	// Cria monitor de recursos
	monitor := monitoring.NewMonitor(time.Second)

	// Canal para receber métricas
	metricsChan := make(chan monitoring.ResourceMetrics, 10)
	err = monitor.Subscribe(metricsChan)
	if err != nil {
		log.Fatalf("Erro ao subscrever no monitor: %v", err)
	}

	// Inicia monitoramento
	err = monitor.Start()
	if err != nil {
		log.Fatalf("Erro ao iniciar monitoramento: %v", err)
	}

	// Canal para capturar sinais de interrupção
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("Monitorando recursos (Ctrl+C para parar)...")

	// Loop principal
	for {
		select {
		case <-sigChan:
			fmt.Println("\nParando monitoramento...")
			monitor.Stop()
			return
		case metrics := <-metricsChan:
			printMetrics(metrics)
		}
	}
}

func printMetrics(metrics monitoring.ResourceMetrics) {
	fmt.Printf("\033[2J\033[H") // Limpa a tela
	fmt.Println("=== Métricas de Sistema ===")
	fmt.Printf("Timestamp: %s\n\n", metrics.Timestamp.Format(time.RFC3339))

	fmt.Println("CPU Usage:")
	for core, usage := range metrics.CPUUsage {
		fmt.Printf("  Core %d: %.2f%%\n", core, usage)
	}

	fmt.Println("\nMemória:")
	fmt.Printf("  Usada: %.2f GB\n", float64(metrics.MemoryUsed)/(1024*1024*1024))
	fmt.Printf("  Disponível: %.2f GB\n", float64(metrics.MemoryAvailable)/(1024*1024*1024))

	if metrics.GPUUsage > 0 {
		fmt.Println("\nGPU:")
		fmt.Printf("  Usage: %.2f%%\n", metrics.GPUUsage)
		fmt.Printf("  Memória Usada: %.2f GB\n", float64(metrics.GPUMemoryUsed)/(1024*1024*1024))
	}
}

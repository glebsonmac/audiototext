package inference

import (
	"github.com/josealecrim/audiototext/internal/models"
)

// ExecutionProvider represents the available execution providers for ONNX Runtime
type ExecutionProvider string

const (
	// CPUExecutionProvider represents the CPU execution provider
	CPUExecutionProvider ExecutionProvider = "CPUExecutionProvider"
	// CUDAExecutionProvider represents the CUDA execution provider for NVIDIA GPUs
	CUDAExecutionProvider ExecutionProvider = "CUDAExecutionProvider"
	// OpenVINOExecutionProvider represents the OpenVINO execution provider for Intel hardware
	OpenVINOExecutionProvider ExecutionProvider = "OpenVINOExecutionProvider"
)

// SessionConfig holds the configuration for an ONNX Runtime session
type SessionConfig struct {
	// ExecutionProvider specifies which hardware to use for inference
	ExecutionProvider ExecutionProvider
	// InterOpNumThreads specifies the number of threads used to parallelize the execution
	// of the graph (across nodes)
	InterOpNumThreads int
	// IntraOpNumThreads specifies the number of threads used to parallelize the execution
	// within nodes
	IntraOpNumThreads int
	// GraphOptimizationLevel specifies the graph optimization level
	GraphOptimizationLevel int
	// EnableMemoryPattern enables memory pattern optimization
	EnableMemoryPattern bool
	// EnableCPUMemArena enables CPU memory arena
	EnableCPUMemArena bool
	// EnableProfiling enables profiling
	EnableProfiling bool
}

// BatchConfig holds the configuration for batch processing
type BatchConfig struct {
	// MaxBatchSize is the maximum number of samples in a batch
	MaxBatchSize int
	// DynamicBatching enables dynamic batching
	DynamicBatching bool
	// MaxLatencyMs is the maximum latency in milliseconds for dynamic batching
	MaxLatencyMs int
	// PreferredBatchSizes is a list of preferred batch sizes
	PreferredBatchSizes []int
}

// Session represents an ONNX Runtime inference session
type Session struct {
	// Model is the ONNX model being used
	Model *models.ONNXModel
	// Config is the session configuration
	Config SessionConfig
	// BatchConfig is the batch processing configuration
	BatchConfig BatchConfig
	// InputShape is the shape of the input tensor
	InputShape []int64
	// OutputShape is the shape of the output tensor
	OutputShape []int64
	// session is the underlying ONNX Runtime session
	session interface{}
}

// Result represents the result of an inference
type Result struct {
	// Transcription is the transcribed text
	Transcription string
	// Confidence is the confidence score (0-1)
	Confidence float32
	// TimestampStart is the start time of the audio segment
	TimestampStart float32
	// TimestampEnd is the end time of the audio segment
	TimestampEnd float32
	// ProcessingTime is the time taken to process the inference
	ProcessingTime float32
}

// Stats represents statistics about the inference session
type Stats struct {
	// TotalInferences is the total number of inferences performed
	TotalInferences int64
	// TotalProcessingTime is the total time spent processing inferences
	TotalProcessingTime float64
	// AverageProcessingTime is the average time per inference
	AverageProcessingTime float64
	// TotalBatches is the total number of batches processed
	TotalBatches int64
	// AverageBatchSize is the average batch size
	AverageBatchSize float64
	// PeakMemoryUsage is the peak memory usage in bytes
	PeakMemoryUsage int64
	// CurrentMemoryUsage is the current memory usage in bytes
	CurrentMemoryUsage int64
}

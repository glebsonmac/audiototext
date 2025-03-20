package helpers

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

// TestConfig holds common test configuration
type TestConfig struct {
	TempDir     string
	ModelDir    string
	CacheDir    string
	TestDataDir string
}

// SetupTestEnv creates a test environment with temporary directories
func SetupTestEnv(t *testing.T) *TestConfig {
	t.Helper()

	tempDir := t.TempDir()
	cfg := &TestConfig{
		TempDir:     tempDir,
		ModelDir:    filepath.Join(tempDir, "models"),
		CacheDir:    filepath.Join(tempDir, "cache"),
		TestDataDir: filepath.Join(tempDir, "testdata"),
	}

	dirs := []string{cfg.ModelDir, cfg.CacheDir, cfg.TestDataDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create directory %s: %v", dir, err)
		}
	}

	return cfg
}

// WithTimeout runs a test with timeout
func WithTimeout(t *testing.T, timeout time.Duration, f func(ctx context.Context)) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		f(ctx)
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		t.Fatal("test timed out")
	}
}

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// AssertEqual fails the test if expected != actual
func AssertEqual(t *testing.T, expected, actual interface{}) {
	t.Helper()
	if expected != actual {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

// GetTestResourcePath returns the absolute path to a test resource
func GetTestResourcePath(t *testing.T, relativePath string) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to get caller info")
	}
	return filepath.Join(filepath.Dir(filename), "..", "resources", relativePath)
}

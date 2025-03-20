package grpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/josealecrim/audiototext/internal/grpc"
	"github.com/josealecrim/audiototext/test/helpers"
)

func setupTestServer(t *testing.T) (grpc.Server, string) {
	// Find available port
	listener, err := net.Listen("tcp", "localhost:0")
	helpers.AssertNoError(t, err)

	// Create server
	server, err := grpc.NewServer(grpc.ServerConfig{
		MaxConcurrentStreams: 10,
		KeepAliveTime:        time.Minute,
		KeepAliveTimeout:     10 * time.Second,
	})
	helpers.AssertNoError(t, err)

	// Start server
	go func() {
		if err := server.Serve(listener); err != nil {
			t.Logf("Server stopped: %v", err)
		}
	}()

	return server, listener.Addr().String()
}

func TestServerStartStop(t *testing.T) {
	server, addr := setupTestServer(t)
	defer server.Stop()

	// Create client
	client, err := grpc.NewClient(addr, grpc.ClientConfig{
		Timeout:          5 * time.Second,
		KeepAliveTime:    time.Minute,
		MaxRetryAttempts: 3,
	})
	helpers.AssertNoError(t, err)
	defer client.Close()

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx)
	helpers.AssertNoError(t, err)
}

func TestAudioStreaming(t *testing.T) {
	server, addr := setupTestServer(t)
	defer server.Stop()

	client, err := grpc.NewClient(addr, grpc.ClientConfig{
		Timeout:          5 * time.Second,
		KeepAliveTime:    time.Minute,
		MaxRetryAttempts: 3,
	})
	helpers.AssertNoError(t, err)
	defer client.Close()

	helpers.WithTimeout(t, 30*time.Second, func(ctx context.Context) {
		// Create audio stream
		stream, err := client.StartStream(ctx)
		helpers.AssertNoError(t, err)

		// Send some audio data
		sampleRate := 16000
		duration := 5 * time.Second
		samples := make([]float32, int(duration.Seconds()*float64(sampleRate)))

		err = stream.Send(samples)
		helpers.AssertNoError(t, err)

		// Get transcription
		result, err := stream.Receive()
		helpers.AssertNoError(t, err)

		if len(result.Text) == 0 {
			t.Error("Transcription result should not be empty")
		}

		// Close stream
		err = stream.Close()
		helpers.AssertNoError(t, err)
	})
}

func TestLoadBalancing(t *testing.T) {
	// Start multiple servers
	numServers := 3
	var servers []grpc.Server
	var addrs []string

	for i := 0; i < numServers; i++ {
		server, addr := setupTestServer(t)
		servers = append(servers, server)
		addrs = append(addrs, addr)
		defer server.Stop()
	}

	// Create load-balanced client
	client, err := grpc.NewLoadBalancedClient(addrs, grpc.ClientConfig{
		Timeout:          5 * time.Second,
		KeepAliveTime:    time.Minute,
		MaxRetryAttempts: 3,
	})
	helpers.AssertNoError(t, err)
	defer client.Close()

	// Test multiple requests
	helpers.WithTimeout(t, time.Minute, func(ctx context.Context) {
		for i := 0; i < 10; i++ {
			stream, err := client.StartStream(ctx)
			helpers.AssertNoError(t, err)

			err = stream.Send(make([]float32, 16000))
			helpers.AssertNoError(t, err)

			_, err = stream.Receive()
			helpers.AssertNoError(t, err)

			err = stream.Close()
			helpers.AssertNoError(t, err)
		}
	})
}

func TestReconnection(t *testing.T) {
	server, addr := setupTestServer(t)

	client, err := grpc.NewClient(addr, grpc.ClientConfig{
		Timeout:          5 * time.Second,
		KeepAliveTime:    time.Minute,
		MaxRetryAttempts: 3,
	})
	helpers.AssertNoError(t, err)
	defer client.Close()

	helpers.WithTimeout(t, time.Minute, func(ctx context.Context) {
		// Test initial connection
		err := client.Ping(ctx)
		helpers.AssertNoError(t, err)

		// Stop server
		server.Stop()

		// Start new server
		newServer, _ := setupTestServer(t)
		defer newServer.Stop()

		// Wait for reconnection
		time.Sleep(2 * time.Second)

		// Test connection again
		err = client.Ping(ctx)
		helpers.AssertNoError(t, err)
	})
}

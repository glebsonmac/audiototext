package grpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// ServerConfig contains configuration for the gRPC server
type ServerConfig struct {
	MaxConcurrentStreams uint32
	KeepAliveTime        time.Duration
	KeepAliveTimeout     time.Duration
}

// ClientConfig contains configuration for the gRPC client
type ClientConfig struct {
	Timeout          time.Duration
	KeepAliveTime    time.Duration
	MaxRetryAttempts int
}

// Server represents a gRPC server
type Server interface {
	Serve(net.Listener) error
	Stop()
}

// Client represents a gRPC client
type Client interface {
	Ping(ctx context.Context) error
	StartStream(ctx context.Context) (Stream, error)
	Close() error
}

// Stream represents an audio streaming session
type Stream interface {
	Send(samples []float32) error
	Receive() (*TranscriptionResult, error)
	Close() error
}

// TranscriptionResult contains the result of audio transcription
type TranscriptionResult struct {
	Text string
}

// server implements the Server interface
type server struct {
	grpcServer *grpc.Server
	config     ServerConfig
}

// client implements the Client interface
type client struct {
	conn   *grpc.ClientConn
	config ClientConfig
}

// stream implements the Stream interface
type stream struct {
	client *client
}

// NewServer creates a new gRPC server
func NewServer(config ServerConfig) (Server, error) {
	s := grpc.NewServer(
		grpc.MaxConcurrentStreams(config.MaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    config.KeepAliveTime,
			Timeout: config.KeepAliveTimeout,
		}),
	)

	return &server{
		grpcServer: s,
		config:     config,
	}, nil
}

// NewClient creates a new gRPC client
func NewClient(addr string, config ClientConfig) (Client, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(config.Timeout),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: config.KeepAliveTime,
		}),
	)
	if err != nil {
		return nil, err
	}
	return &client{conn: conn, config: config}, nil
}

// NewLoadBalancedClient creates a new load-balanced gRPC client
func NewLoadBalancedClient(addrs []string, config ClientConfig) (Client, error) {
	// TODO: Implement load balancing
	return NewClient(addrs[0], config)
}

// Implementation of server methods
func (s *server) Serve(listener net.Listener) error {
	return s.grpcServer.Serve(listener)
}

func (s *server) Stop() {
	s.grpcServer.GracefulStop()
}

// Implementation of client methods
func (c *client) Ping(ctx context.Context) error {
	// TODO: Implement ping
	return nil
}

func (c *client) StartStream(ctx context.Context) (Stream, error) {
	return &stream{client: c}, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

// Implementation of stream methods
func (s *stream) Send(samples []float32) error {
	// TODO: Implement send
	return nil
}

func (s *stream) Receive() (*TranscriptionResult, error) {
	// TODO: Implement receive
	return &TranscriptionResult{
		Text: "Sample transcription",
	}, nil
}

func (s *stream) Close() error {
	// TODO: Implement close
	return nil
}

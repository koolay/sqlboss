package agent

import (
	"fmt"
	"time"

	natsd "github.com/nats-io/nats-server/v2/server"
	stand "github.com/nats-io/nats-streaming-server/server"
)

type StreamingServerConfig struct {
	ClusterID string
	Name      string
	Host      string
	Port      int
	Debug     bool
	MaxConn   int
}

func DefaultStreamingServerConfig() *StreamingServerConfig {
	return &StreamingServerConfig{
		ClusterID: "local",
		Name:      "defaultServer",
		Host:      "127.0.0.1",
		Port:      4222,
		Debug:     false,
	}
}

type StreamingServer struct {
	serverCfg *StreamingServerConfig
}

func NewStreamingServer(serverCfg *StreamingServerConfig) *StreamingServer {
	return &StreamingServer{
		serverCfg: serverCfg,
	}
}

func (s *StreamingServer) Start() error {
	fmt.Printf("nats-streaming-server version %s, ", stand.VERSION)
	sOpts := stand.GetDefaultOptions()
	// Force the streaming server to setup its own signal handler
	sOpts.HandleSignals = true
	// // override the NoSigs for NATS since Streaming has its own signal handler
	// nOpts.NoSigs = true
	// Without this option set to true, the logger is not configured.
	sOpts.EnableLogging = true
	sOpts.ClientHBInterval = 3 * time.Second
	sOpts.ClientHBTimeout = 60 * time.Second
	sOpts.ID = s.serverCfg.ClusterID
	nOpts := &natsd.Options{
		ServerName: s.serverCfg.Name,
		Host:       s.serverCfg.Host,
		Port:       s.serverCfg.Port,
		Debug:      s.serverCfg.Debug,
		MaxConn:    s.serverCfg.MaxConn,
	}
	// This will invoke RunServerWithOpts but on Windows, may run it as a service.
	if _, err := stand.Run(sOpts, nOpts); err != nil {
		return err
	}

	return nil
}

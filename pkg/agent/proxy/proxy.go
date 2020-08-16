package proxy

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/pingcap/errors"
	"github.com/siddontang/go-mysql/server"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	cfg          Config
	logger       *logrus.Logger
	ser          *ProxyServer
	handler      server.Handler
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock *sync.Mutex
}

func NewProxy(cfg Config, logger *logrus.Logger, ser *ProxyServer, handler server.Handler) (*Proxy, error) {
	return &Proxy{
		cfg:          cfg,
		logger:       logger,
		ser:          ser,
		handler:      handler,
		shutdownLock: &sync.Mutex{},
		shutdownCh:   make(chan struct{}),
	}, nil
}

func (p *Proxy) Shutdown() error {
	p.logger.Info("shutting down db proxy")
	p.shutdownLock.Lock()
	defer p.shutdownLock.Unlock()

	if p.shutdown {
		return nil
	}

	p.shutdown = true
	close(p.shutdownCh)
	return nil
}

func (p *Proxy) Start() error {
	listener, err := net.Listen("tcp", p.cfg.Addr)
	if err != nil {
		return errors.Wrapf(err, "failed to listen: %s", p.cfg.Addr)
	}

	for {
		fmt.Println("new accept")
		conn, err := listener.Accept()
		if err != nil {
			if p.shutdown {
				return nil
			}

			p.logger.WithError(err).Error("failed to accept conn", err)
			continue
		}

		go p.handleConn(conn)
	}
}

func (p *Proxy) handleConn(conn net.Conn) {
	log.Println("handleConn")
	dbconn, err := p.ser.NewConnect(conn, p.handler)
	if err != nil {
		log.Printf("Connection error: %v", err)
		return
	}

	defer func() {
		if dbconn.Conn != nil {
			dbconn.Close()
		}
	}()

	for {
		select {
		case <-p.shutdownCh:
			return
		default:
		}

		if dbconn.Conn == nil {
			p.logger.Debug("closed connection")
			return
		}

		if err := dbconn.HandleCommand(); err != nil {
			p.logger.WithError(err).Error("failed to handle command")
		}
	}
}

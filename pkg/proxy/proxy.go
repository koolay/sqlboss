package proxy

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/koolay/sqlboss/pkg/conf"
	"github.com/koolay/sqlboss/pkg/message"
	"github.com/pingcap/errors"
	"github.com/siddontang/go-mysql/server"
	"github.com/sirupsen/logrus"
)

type Proxy struct {
	cfg          *conf.Config
	mysqlCfg     *MysqlServerConfig
	logger       *logrus.Logger
	ser          *mysqlSession
	handler      server.Handler
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock *sync.Mutex
}

func NewProxy(cfg *conf.Config,
	logger *logrus.Logger,
	mysqlCfg *MysqlServerConfig,
	eventBus message.EventBus,
) (*Proxy, error) {

	sess := newMysqlSession(mysqlCfg.Version, mysqlCfg)
	handler, err := newMysqlHandler(cfg,
		Connection{
			Host:     cfg.DB.Host,
			Port:     cfg.DB.Port,
			User:     cfg.DB.User,
			Password: cfg.DB.Password,
			Database: cfg.DB.Database,
		}, eventBus, logger)
	if err != nil {
		return nil, err
	}

	return &Proxy{
		cfg:          cfg,
		mysqlCfg:     mysqlCfg,
		logger:       logger,
		ser:          sess,
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
	log.Println("start listen", p.mysqlCfg.Addr)
	listener, err := net.Listen("tcp", p.mysqlCfg.Addr)
	if err != nil {
		return errors.Wrapf(err, "failed to listen: %s", p.mysqlCfg.Addr)
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
	dbconn, err := p.ser.newConnect(conn, p.handler)
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

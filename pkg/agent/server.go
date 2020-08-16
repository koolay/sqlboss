package agent

import (
	"log"

	"github.com/koolay/sqlboss/pkg/agent/proxy"
	"github.com/koolay/sqlboss/pkg/conf"
	"github.com/sirupsen/logrus"
)

type Server struct {
	proxys *proxy.ProxyServer
	cfg    *conf.Config
	logger *logrus.Logger
}

func NewServer(cfg *conf.Config, mysqlVersion string) *Server {
	return &Server{
		cfg: cfg,
		proxys: proxy.NewProxyServer(mysqlVersion, proxy.Config{
			Addr:     "127.0.0.1:3309",
			User:     "root",
			Password: "123",
			TargetConnection: proxy.Connection{
				Host:     cfg.DB.Host,
				Port:     cfg.DB.Port,
				User:     cfg.DB.User,
				Password: cfg.DB.Password,
				Database: cfg.DB.Database,
			},
		}),
	}
}

func (s *Server) Start() {
	handler, err := proxy.NewMysqlHandler(proxy.Connection{}, s.logger)
	if err != nil {
		s.logger.WithError(err).Error("failed to new handler")
		return
	}

	proxier, err := proxy.NewProxy(proxy.Config{}, s.logger, s.proxys, handler)
	if err != nil {
		s.logger.WithError(err).Error("failed to new proxy")
		return
	}

	go func() {
		log.Println("start db proxy")
		if err := proxier.Start(); err != nil {
			s.logger.WithError(err).Error("failed to start proxy")
		}
	}()

}

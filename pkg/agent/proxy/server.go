package proxy

import (
	"net"
	"time"

	"crypto/tls"

	"github.com/pingcap/errors"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/server"
	"github.com/siddontang/go-mysql/test_util/test_keys"
)

var (
	tlsConf = server.NewServerTLSConfig(test_keys.CaPem,
		test_keys.CertPem,
		test_keys.KeyPem,
		tls.VerifyClientCertIfGiven)
)

type RemoteThrottleProvider struct {
	*server.InMemoryProvider
	delay int // in milliseconds
}

func (m *RemoteThrottleProvider) GetCredential(username string) (password string, found bool, err error) {
	time.Sleep(time.Millisecond * time.Duration(m.delay))
	return m.InMemoryProvider.GetCredential(username)
}

type ProxyServer struct {
	s   *server.Server
	cfg Config
}

func NewProxyServer(version string, cfg Config) *ProxyServer {
	return &ProxyServer{
		cfg: cfg,
		s: server.NewServer(version,
			mysql.DEFAULT_COLLATION_ID,
			mysql.AUTH_CACHING_SHA2_PASSWORD,
			test_keys.PubPem,
			tlsConf),
	}

}

func (s *ProxyServer) NewConnect(conn net.Conn, handler server.Handler) (*server.Conn, error) {
	ser := server.NewServer("5.7.28",
		mysql.DEFAULT_COLLATION_ID,
		mysql.AUTH_CACHING_SHA2_PASSWORD,
		test_keys.PubPem,
		tlsConf)

	credentialProvider := &RemoteThrottleProvider{server.NewInMemoryProvider(), 10 + 50}
	credentialProvider.AddUser(s.cfg.User, s.cfg.Password)
	dbconn, err := server.NewCustomizedConn(conn, ser, credentialProvider, handler)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new customizedConnection")
	}

	return dbconn, nil
}

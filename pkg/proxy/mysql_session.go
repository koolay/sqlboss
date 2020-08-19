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

type mysqlSession struct {
	s        *server.Server
	mysqlCfg *MysqlServerConfig
}

func newMysqlSession(version string, mysqlCfg *MysqlServerConfig) *mysqlSession {
	return &mysqlSession{
		mysqlCfg: mysqlCfg,
		s: server.NewServer(version,
			mysql.DEFAULT_COLLATION_ID,
			mysql.AUTH_CACHING_SHA2_PASSWORD,
			test_keys.PubPem,
			tlsConf),
	}

}

func (s *mysqlSession) newConnect(conn net.Conn, handler server.Handler) (*server.Conn, error) {
	ser := server.NewServer("5.7.28",
		mysql.DEFAULT_COLLATION_ID,
		mysql.AUTH_CACHING_SHA2_PASSWORD,
		test_keys.PubPem,
		tlsConf)

	credentialProvider := &RemoteThrottleProvider{server.NewInMemoryProvider(), 10 + 50}
	credentialProvider.AddUser(s.mysqlCfg.User, s.mysqlCfg.Password)
	dbconn, err := server.NewCustomizedConn(conn, ser, credentialProvider, handler)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new customizedConnection")
	}

	return dbconn, nil
}

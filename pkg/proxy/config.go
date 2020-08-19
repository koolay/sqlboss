package proxy

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type MysqlServerConfig struct {
	Version          string
	Addr             string
	User             string
	Password         string
	TargetConnection Connection
}

func DefaultConfig() MysqlServerConfig {
	return MysqlServerConfig{
		Addr:     "0.0.0.0:3306",
		User:     "root",
		Password: "dev",
	}
}

package proxy

type Connection struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Config struct {
	Addr             string
	User             string
	Password         string
	TargetConnection Connection
}

func DefaultConfig() Config {
	return Config{
		Addr:     "0.0.0.0:3306",
		User:     "root",
		Password: "dev",
	}
}

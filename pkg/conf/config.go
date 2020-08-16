package conf

import (
	"fmt"
	"log"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
	"gopkg.in/ini.v1"
)

const (
	defaultEnvCode         = "dev"
	RuntimeDev     Runtime = "dev"
	RuntimeProd    Runtime = "prod"

	DefaultNatsURL     = "nats://127.0.0.1:4222"
	DefaultSqlTopic    = "sqlcmd"
	DefaultNatsCluster = "defaultCluster"
)

type Runtime string
type BoolString string

func (bs BoolString) ToBool() bool {
	bv := string(bs)
	return bv == "true" || bv == "TRUE" || bv == "1" || bv == "True"
}

type App struct {
	Name string `ini:"name" json:"name"`
	// Runtime dev, test, release
	Runtime Runtime `ini:"runtime" json:"runtime"`
}

type Log struct {
	Handler   string `ini:"handler" json:"handler"`
	Level     string `ini:"level" json:"level"`
	SentryDSN string `ini:"sentry_dsn" json:"sentry_dsn"`
}

type DB struct {
	Database string `ini:"database" json:"database"`
	Host     string `ini:"host" json:"host"`
	Port     int    `ini:"port" json:"port"`
	Password string `ini:"password" json:"password"`
	User     string `ini:"user" json:"user"`

	MaxIdleConns int `ini:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns int `ini:"max_open_conns" json:"max_open_conns"`
}

type Nats struct {
	URL       string `ini:"url"`
	ClusterID string `ini:"cluster_id"`
}

type Stream struct {
	Topic string `ini:"topic"`
}

type Config struct {
	App App `ini:"App"`
	Log struct {
		Handler   string `ini:"handler"`
		Level     string `ini:"level"`
		SentryDSN string `ini:"sentry_dsn"`
	} `ini:"Log"`
	DB     DB     `ini:"DB"`
	Nats   Nats   `ini:"Nats"`
	Stream Stream `ini:"Stream"`
}

func (d *DB) Default() DB {
	return DB{
		Host:         "localhost",
		Port:         3306,
		Database:     "demo",
		User:         "root",
		Password:     "dev",
		MaxIdleConns: 25,
		MaxOpenConns: 50,
	}
}

func (d *DB) DSN() string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s", d.User, d.Password, d.Host, d.Port, d.Database)
}

// LoadConfig 读取统一的配置文件,app.config
func LoadConfig(cfgFile string) (*Config, error) {
	var cfg Config

	f, err := ini.LoadSources(ini.LoadOptions{
		// 忽略行内注释字符:#
		IgnoreInlineComment: true,
	}, cfgFile)
	if err != nil {
		return nil, err
	}

	if err = f.MapTo(&cfg); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to bind ini config, file: %s", cfgFile))
	}

	defaultCfg := defaultConfiguration()
	if err = mergo.Merge(&cfg, defaultCfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func InitDeafultCfgFile(filePath string) error {
	appCfg := defaultConfiguration()
	cfg := ini.Empty()
	err := ini.ReflectFrom(cfg, appCfg)
	if err != nil {
		return err
	}
	log.Printf("init configration file: %s\n", filePath)
	return cfg.SaveTo(filePath)
}

func defaultConfiguration() *Config {
	return &Config{
		Log: struct {
			Handler   string `ini:"handler"`
			Level     string `ini:"level"`
			SentryDSN string `ini:"sentry_dsn"`
		}{
			Handler: "console",
			Level:   "INFO",
		},
		App: App{
			Name:    "sqlboss",
			Runtime: RuntimeProd,
		},
		DB: DB{
			Database:     "demo",
			Host:         "localhost",
			Port:         3306,
			User:         "root",
			Password:     "dev",
			MaxOpenConns: 100,
			MaxIdleConns: 0,
		},
		Nats: Nats{
			URL:       DefaultNatsURL,
			ClusterID: DefaultNatsCluster,
		},
		Stream: Stream{
			Topic: DefaultSqlTopic,
		},
	}
}

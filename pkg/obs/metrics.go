package obs

import (
	"fmt"
	"net/http"

	"github.com/pingcap/errors"

	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	// metricstdout "go.opentelemetry.io/otel/exporters/metric/stdout"
)

const (
	defaultExporterPath = "/metrics"
	defaultExporterPort = 2222
)

type Config struct {
	ExporterPath string
	ExportPort   int
}

type Telemetry struct {
	// pusher   *push.Controller
	cfg      *Config
	Meter    metric.Meter
	exporter *prometheus.Exporter
}

func NewTelemetry(cfg Config) (*Telemetry, error) {
	exporter, err := prometheus.InstallNewPipeline(prometheus.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize prometheus exporter")
	}

	return &Telemetry{
		cfg:      &cfg,
		Meter:    global.Meter("sample"),
		exporter: exporter,
	}, nil
}

func (t *Telemetry) Serve() error {
	exporterPath := defaultExporterPath
	if t.cfg.ExporterPath != "" {
		exporterPath = t.cfg.ExporterPath
	}

	http.HandleFunc(exporterPath, t.exporter.ServeHTTP)
	addr := fmt.Sprintf(":%d", defaultExporterPort)
	if t.cfg.ExportPort == 0 {
		addr = fmt.Sprintf(":%d", t.cfg.ExportPort)
	}

	fmt.Println("Prometheus server running on", addr)
	return http.ListenAndServe(addr, nil)
}

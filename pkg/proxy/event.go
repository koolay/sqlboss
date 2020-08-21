package proxy

import (
	"context"

	"github.com/koolay/sqlboss/pkg/lineage"
	"github.com/koolay/sqlboss/pkg/message"
	"github.com/koolay/sqlboss/pkg/parse/fingerprint"
	"github.com/koolay/sqlboss/pkg/store"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/api/kv"
	"go.opentelemetry.io/otel/api/metric"
)

// ParseOnSQLEventHandler sql命令事件处理, 血缘分析
type ParseOnSQLEventHandler struct {
	logger     *logrus.Entry
	commandBus message.CommandBus

	sqlCounter metric.Int64Counter
	sqlLatency metric.Float64ValueRecorder
	// storager Storager
}

func NewParseOnSQLEventHandler(logger *logrus.Entry, commandBus message.CommandBus, meter metric.Meter) ParseOnSQLEventHandler {
	metric.Must(meter).NewInt64Counter("sql_request_count")

	return ParseOnSQLEventHandler{
		commandBus: commandBus,
		logger:     logger.WithField("name", "ParseOnSQLEventHandler"),
		sqlCounter: metric.Must(meter).NewInt64Counter("sql_request_count"),
		sqlLatency: metric.Must(meter).NewFloat64ValueRecorder("sql_request_latency"),
	}
}

func (s ParseOnSQLEventHandler) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "ParseOnSQLEventHandler"
}

func (ParseOnSQLEventHandler) NewEvent() interface{} {
	return &SqlCommand{}
}

func (s ParseOnSQLEventHandler) Handle(ctx context.Context, event interface{}) error {
	data := event.(*SqlCommand)
	// 解析 event.SQL, 1. 生成血缘RDF数据 2.LogCommand
	sqlFingerprint := fingerprint.Fingerprint(data.SQL)
	s.logger.WithFields(logrus.Fields{
		"sql":             data.SQL,
		"sql_fingerprint": sqlFingerprint,
	}).Info("received sql command")

	lables := kv.Any("sql", sqlFingerprint)

	s.sqlCounter.Add(ctx, 1, lables)
	s.sqlLatency.Record(ctx, float64(data.Duration), lables)

	// 1. write log
	logMsg := &store.LogCommand{
		Env:              data.Env,
		App:              data.App,
		Database:         data.Database,
		SQL:              data.SQL,
		SqlFingerprint:   fingerprint.Fingerprint(data.SQL),
		User:             data.User,
		Duration:         data.Duration,
		Occtime:          data.Occtime,
		Tables:           []string{},
		Fields:           []string{},
		PerformanceScore: 0.0,
	}

	if err := s.commandBus.Send(ctx, logMsg); err != nil {
		return err
	}

	// 2. 生成血缘关系
	quads := &lineage.Quads{}
	return s.commandBus.Send(ctx, quads)
}

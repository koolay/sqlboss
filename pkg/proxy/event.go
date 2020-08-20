package proxy

import (
	"context"

	"github.com/koolay/sqlboss/pkg/lineage"
	"github.com/koolay/sqlboss/pkg/message"
	"github.com/koolay/sqlboss/pkg/parse/fingerprint"
	"github.com/koolay/sqlboss/pkg/store"
	"github.com/sirupsen/logrus"
)

// ParseOnSQLEventHandler sql命令事件处理, 血缘分析
type ParseOnSQLEventHandler struct {
	logger     *logrus.Entry
	commandBus message.CommandBus
	// storager Storager
}

func NewParseOnSQLEventHandler(logger *logrus.Entry, commandBus message.CommandBus) ParseOnSQLEventHandler {
	return ParseOnSQLEventHandler{
		commandBus: commandBus,
		logger:     logger.WithField("name", "ParseOnSQLEventHandler"),
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
	s.logger.WithField("sql", data.SQL).Info("received sql command")

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
		Table:            "",
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

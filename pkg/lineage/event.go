package lineage

import (
	"context"

	"github.com/koolay/sqlboss/pkg/proto"
	"github.com/sirupsen/logrus"
)

// LineageOnSQLEventHandler sql命令事件处理, 血缘分析
type LineageOnSQLEventHandler struct {
	logger *logrus.Entry
	// storager Storager
}

func NewLineageOnSQLEventHandler(logger *logrus.Entry) LineageOnSQLEventHandler {
	return LineageOnSQLEventHandler{
		logger: logger.WithField("name", "LineageOnSQLEventHandler"),
	}
}

func (s LineageOnSQLEventHandler) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "LineageOnSQLEventHandler"
}

func (LineageOnSQLEventHandler) NewEvent() interface{} {
	return &proto.SqlCommand{}
}

func (s LineageOnSQLEventHandler) Handle(ctx context.Context, event interface{}) error {
	data := event.(*proto.SqlCommand)
	// 解析 event.SQL, 生成血缘RDF数据
	s.logger.WithField("sql", data.SQL).Info("received sql command")
	return nil
}

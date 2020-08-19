package lineage

import (
	"context"

	"github.com/sirupsen/logrus"
)

// LineageCommandHandler 血缘命令处理, 血缘分析
type LineageCommandHandler struct {
	logger *logrus.Entry
	// storager Storager
}

func NewLineageCommandHandler(logger *logrus.Entry) LineageCommandHandler {
	return LineageCommandHandler{
		logger: logger.WithField("name", "LineageCommandHandler"),
	}
}

func (s LineageCommandHandler) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "LineageCommandHandler"
}

func (LineageCommandHandler) NewCommand() interface{} {
	return &Quads{}
}

func (s LineageCommandHandler) Handle(ctx context.Context, event interface{}) error {
	data := event.(*Quads)
	// 解析 event.SQL, 生成血缘RDF数据
	s.logger.WithField("sql", data).Info("received lineage command")
	return nil
}

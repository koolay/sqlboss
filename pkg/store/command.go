package store

import (
	"context"

	"github.com/sirupsen/logrus"
)

// StoreCommandHandler store命令处理, 日志存储
type StoreCommandHandler struct {
	logger *logrus.Entry
	// storager Storager
}

func NewStoreCommandHandler(logger *logrus.Entry) StoreCommandHandler {
	return StoreCommandHandler{
		logger: logger.WithField("name", "StoreCommandHandler"),
	}
}

func (s StoreCommandHandler) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "StoreCommandHandler"
}

func (StoreCommandHandler) NewCommand() interface{} {
	return &LogCommand{}
}

func (s StoreCommandHandler) Handle(ctx context.Context, cmd interface{}) error {
	data := cmd.(*LogCommand)
	// 解析 event.SQL, 用soar分析sql性能, 生成message对象, 并存储
	// s.storager.Insert()
	s.logger.WithField("sql", data.SQL).Info("received sql command")
	return nil
}

package store

import (
	"context"

	"github.com/koolay/sqlboss/pkg/proto"
	"github.com/sirupsen/logrus"
)

// StoreOnSQLEventHandler sql命令事件处理, 日志存储
type StoreOnSQLEventHandler struct {
	logger *logrus.Entry
	// storager Storager
}

func NewStoreOnSQLEventHandler(logger *logrus.Entry) StoreOnSQLEventHandler {
	return StoreOnSQLEventHandler{
		logger: logger.WithField("name", "StoreOnSQLEventHandler"),
	}
}

func (s StoreOnSQLEventHandler) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "StoreOnSQLEventHandler"
}

func (StoreOnSQLEventHandler) NewEvent() interface{} {
	return &proto.SqlCommand{}
}

func (s StoreOnSQLEventHandler) Handle(ctx context.Context, event interface{}) error {
	data := event.(*proto.SqlCommand)
	// 解析 event.SQL, 用soar分析sql性能, 生成message对象, 并存储
	// s.storager.Insert()
	s.logger.WithField("sql", data.SQL).Info("received sql command")
	return nil
}

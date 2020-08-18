package lineage

import (
	"context"
	"fmt"

	"github.com/koolay/sqlboss/pkg/proto"
)

// LineageOnSQLEventHandler sql命令事件处理, 血缘分析
type LineageOnSQLEventHandler struct {
	// storager Storager
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
	fmt.Println("received sql command")
	fmt.Println("LineageOnSQLEventHandler", data)
	return nil
}

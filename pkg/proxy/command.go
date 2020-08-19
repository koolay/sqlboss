package proxy

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/koolay/sqlboss/pkg/proto"
)

// SQLCommandHandler 处理sql命令事件
type SQLCommandHandler struct {
	eventBus *cqrs.EventBus
	// sql query -> send event
}

func NewSQLCommandHandler(eventBus *cqrs.EventBus) SQLCommandHandler {
	return SQLCommandHandler{eventBus: eventBus}
}

func (b SQLCommandHandler) HandlerName() string {
	return "SQLCommandHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b SQLCommandHandler) NewCommand() interface{} {
	return &proto.SqlCommand{}
}

func (b SQLCommandHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*proto.SqlCommand)
	return b.eventBus.Publish(ctx, cmd)
}

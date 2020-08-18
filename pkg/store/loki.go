package store

import "github.com/koolay/sqlboss/pkg/proto"

type Storager interface {
	Insert(data *proto.LogMessage) error
}

type LokiStorager struct {
}

func (lk LokiStorager) Insert(data *proto.LogMessage) error {
	return nil
}

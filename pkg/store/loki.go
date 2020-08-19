package store

type Storager interface {
	Insert(data *LogCommand) error
}

type LokiStorager struct {
}

func (lk LokiStorager) Insert(data *LogCommand) error {
	return nil
}

package proxy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/koolay/sqlboss/pkg/conf"
	"github.com/koolay/sqlboss/pkg/message"
	"github.com/koolay/sqlboss/pkg/proto"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/siddontang/go-mysql/client"
	"github.com/siddontang/go-mysql/mysql"
)

type MysqlHandler struct {
	cfg        *conf.Config
	logger     *logrus.Logger
	conn       *client.Conn
	statements map[int64]*client.Stmt
	commandBus message.CommandBus
}

func newMysqlHandler(cfg *conf.Config,
	targetConn Connection,
	commandBus message.CommandBus,
	logger *logrus.Logger) (*MysqlHandler, error) {
	conn, err := client.Connect(fmt.Sprintf("%s:%d", targetConn.Host, targetConn.Port),
		targetConn.User,
		targetConn.Password,
		targetConn.Database)

	if err != nil {
		return nil, errors.Wrap(err, "failed to connect target")
	}

	return &MysqlHandler{
		logger:     logger,
		commandBus: commandBus,
		conn:       conn,
		cfg:        cfg,
		statements: make(map[int64]*client.Stmt),
	}, nil
}

func (h *MysqlHandler) UseDB(dbName string) error {
	return h.conn.UseDB(dbName)
}

func (h *MysqlHandler) HandleQuery(query string) (*mysql.Result, error) {
	st := time.Now()
	result, err := h.conn.Execute(query)
	data := &proto.SqlCommand{
		App:      h.cfg.App.Name,
		Database: h.cfg.DB.Database,
		User:     h.cfg.DB.User,
		Env:      "dev",
		SQL:      query,
		Duration: time.Since(st).Milliseconds(),
		Occtime:  st.Unix(),
	}

	if err := h.commandBus.Send(context.Background(), data); err != nil {
		// ignore error and continue
		h.logger.WithError(err).Error("failed to send commandBus")
	}

	if err != nil {
		return nil, errors.Wrapf(err, "failed to handleQuery:%s", query)
	}

	return result, nil
}

func (h *MysqlHandler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	h.logger.Debug("HandleFieldList")
	fields, err := h.conn.FieldList(table, fieldWildcard)
	if err != nil {
		return nil, errors.Wrap(err, "failed to HandleFieldList")
	}

	return fields, nil
}

func (h *MysqlHandler) HandleStmtPrepare(query string) (int, int, interface{}, error) {
	h.logger.Debug("HandleStmtPrepare")
	stmt, err := h.conn.Prepare(query)

	if err != nil {
		return 0, 0, nil, errors.Wrap(err, "failed to HandleStmtPrepare")
	}

	id := time.Now().UnixNano()
	h.statements[id] = stmt

	paramNum := stmt.ParamNum()
	colNum := stmt.ColumnNum()

	return paramNum, colNum, id, nil
}

func (h *MysqlHandler) HandleStmtExecute(context interface{}, query string, args []interface{}) (*mysql.Result, error) {
	h.logger.Debug("HandleStmtExecute")
	intContext, ok := context.(int64)
	if !ok {
		log.Printf("Invalid context: %+v", context)
		return nil, fmt.Errorf("Invalid context")
	}

	if stmt, ok := h.statements[intContext]; !ok {
		log.Println("Creating statement on-the-fly and execute it")
		inlineStmt, err := h.conn.Prepare(query)

		if err != nil {
			return nil, errors.Wrap(err, "failed to prepare")
		}

		return inlineStmt.Execute(args...)
	} else {
		return stmt.Execute(args...)
	}
}

func (h *MysqlHandler) HandleStmtClose(context interface{}) error {
	h.logger.Debug("HandleStmtClose")
	intContext, ok := context.(int64)
	if !ok {
		return nil
	}

	stmt, ok := h.statements[intContext]
	if !ok {
		return nil
	}

	if err := stmt.Close(); err != nil {
		return errors.Wrap(err, "failed to HandleStmtClose")
	}

	return nil
}

func (h *MysqlHandler) HandleOtherCommand(cmd byte, data []byte) error {
	h.logger.Debug("HandleOtherCommand")
	return mysql.NewError(
		mysql.ER_UNKNOWN_ERROR,
		fmt.Sprintf("command %d is not supported now", cmd),
	)
}

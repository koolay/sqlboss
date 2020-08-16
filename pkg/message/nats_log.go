package message

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/sirupsen/logrus"
)

type NatsLog struct {
	logger *logrus.Entry
}

func (l *NatsLog) Error(msg string, err error, fields watermill.LogFields) {
	l.logger.WithError(err).WithFields(map[string]interface{}(fields)).Error(msg)
}

func (l *NatsLog) Info(msg string, fields watermill.LogFields) {
	l.logger.WithFields(map[string]interface{}(fields)).Info(msg)
}

func (l *NatsLog) Debug(msg string, fields watermill.LogFields) {
	l.logger.WithFields(map[string]interface{}(fields)).Debug(msg)
}

func (l *NatsLog) Trace(msg string, fields watermill.LogFields) {
	l.logger.WithFields(map[string]interface{}(fields)).Trace(msg)
}

func (l *NatsLog) With(fields watermill.LogFields) watermill.LoggerAdapter {
	l.logger = l.logger.WithFields(map[string]interface{}(fields))
	return l
}

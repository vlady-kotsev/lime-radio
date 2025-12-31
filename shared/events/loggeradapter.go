package events

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

type ZapLoggerAdapter struct {
	logger *zap.Logger
}

func NewZapLoggerAdapter(logger *zap.Logger) watermill.LoggerAdapter {
	return &ZapLoggerAdapter{logger: logger}
}

func (z *ZapLoggerAdapter) Error(msg string, err error, fields watermill.LogFields) {
	zapFields := z.watermillFieldsToZap(fields)
	z.logger.Error(msg, append(zapFields, zap.Error(err))...)
}

func (z *ZapLoggerAdapter) Info(msg string, fields watermill.LogFields) {
	zapFields := z.watermillFieldsToZap(fields)
	z.logger.Info(msg, zapFields...)
}

func (z *ZapLoggerAdapter) Debug(msg string, fields watermill.LogFields) {
	zapFields := z.watermillFieldsToZap(fields)
	z.logger.Debug(msg, zapFields...)
}

func (z *ZapLoggerAdapter) Trace(msg string, fields watermill.LogFields) {
	zapFields := z.watermillFieldsToZap(fields)
	z.logger.Debug(msg, zapFields...) // Zap doesn't have Trace, use Debug
}

func (z *ZapLoggerAdapter) With(fields watermill.LogFields) watermill.LoggerAdapter {
	zapFields := z.watermillFieldsToZap(fields)
	return &ZapLoggerAdapter{
		logger: z.logger.With(zapFields...),
	}
}

func (z *ZapLoggerAdapter) watermillFieldsToZap(fields watermill.LogFields) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}
	return zapFields
}

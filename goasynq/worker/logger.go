package worker

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.uber.org/zap"
)

type TaskLog struct {
	//*zap.Logger
	//zerolog.Logger
	asynq.Logger
}

func NewTaskLog() asynq.Logger {
	//logger, _ := zap.NewProduction()
	////defer logger.Sync()
	//defer func(logger *zap.Logger) {
	//	err := logger.Sync()
	//	if err != nil {
	//		log.Fatal("new task log failed: ", err)
	//	}
	//}(logger)
	//
	//return &TaskLog{
	//	logger,
	//}
	return &TaskLog{}
}

//type TaskLogger interface {
//	Info(args ...interface{})
//	Warn(args ...interface{})
//	Error(args ...interface{})
//	Fatal(args ...interface{})
//	Debug(args ...interface{})
//}

func (l *TaskLog) Print(level zerolog.Level, args ...interface{}) {
	log.WithLevel(level).Msg(fmt.Sprint(args...))
}

func (l *TaskLog) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (l *TaskLog) Info(args ...interface{}) {
	l.Info(zap.InfoLevel, args)
}

func (l *TaskLog) Warn(args ...interface{}) {
	l.Warn(zap.WarnLevel, args)
}

func (l *TaskLog) Error(args ...interface{}) {
	l.Error(zap.ErrorLevel, args)
}

func (l *TaskLog) Fatal(args ...interface{}) {
	l.Fatal(zap.FatalLevel, args)
}

func (l *TaskLog) Debug(args ...interface{}) {
	l.Debug(zap.DebugLevel, args)
}

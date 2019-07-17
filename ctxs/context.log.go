package ctxs

import (
	"github.com/champly/hercules/configs"
	"github.com/sirupsen/logrus"
)

type ILog interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Infoln(args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Debugln(args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Warnln(args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Errorln(args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
}

type Logger struct {
	l *logrus.Logger
}

func NewLogger() *Logger {
	return &Logger{
		l: newLogger(),
	}
}

func newLogger() *logrus.Logger {
	logger := &logrus.Logger{
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		},
		Hooks: make(logrus.LevelHooks),
	}
	if configs.LoggerInfo.Debug {
		logger = &logrus.Logger{
			Formatter: &Formatter{},
			Hooks:     make(logrus.LevelHooks),
		}
	}
	logger.AddHook(NewFileSplitHook())

	lev, err := logrus.ParseLevel(configs.LoggerInfo.Level)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(lev)
	return logger
}

func (l *Logger) Info(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Info(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Infof(format, args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Infoln(args...)
}

func (l *Logger) Debug(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Debug(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Debugf(format, args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Debugln(args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Warn(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Warnf(format, args...)
}

func (l *Logger) Warnln(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Warnln(args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Error(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Errorf(format, args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Errorln(args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Fatal(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Fatalf(format, args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.l.WithFields(logrus.Fields{}).Fatalln(args...)
}

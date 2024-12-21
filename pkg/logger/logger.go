package logger

import (
	"log"
)

type Logger interface {
	Info(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
}

type logger struct {
	prefix string
}

func NewLogger(level string) Logger {
	return &logger{prefix: "[P2PChat] "}
}

func (l *logger) Info(v ...interface{}) {
	log.Println(append([]interface{}{l.prefix, "INFO:"}, v...)...)
}

func (l *logger) Error(v ...interface{}) {
	log.Println(append([]interface{}{l.prefix, "ERROR:"}, v...)...)
}

func (l *logger) Fatal(v ...interface{}) {
	log.Fatal(append([]interface{}{l.prefix, "FATAL:"}, v...)...)
}

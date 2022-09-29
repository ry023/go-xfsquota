package xfsquota

import (
	"fmt"
	"log"
)

type Logger interface {
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
}

type StandardLog struct {
}

func (l *StandardLog) Infof(s string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[Info] %s", s), args...)
}

func (l *StandardLog) Warnf(s string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[Warning] %s", s), args...)
}

func (l *StandardLog) Errorf(s string, args ...interface{}) {
	log.Printf(fmt.Sprintf("[Error] %s", s), args...)
}

package logger

import (
	"io"
	"log"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func New(outInfo, outError io.Writer) *Logger {
	return &Logger{
		Info:  log.New(outInfo, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(outError, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

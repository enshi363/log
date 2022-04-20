package log

import (
	"io"
	"log"
	"log/syslog"
	"os"
)

type Logger log.Logger

var (
	Fatal     = log.Fatal
	Fatalf    = log.Fatalf
	Fatalln   = log.Fatalln
	Flags     = log.Flags
	Output    = log.Output
	Panic     = log.Panic
	Panicf    = log.Panicf
	Panicln   = log.Panicln
	Prefix    = log.Prefix
	Print     = log.Print
	Printf    = log.Printf
	Println   = log.Println
	SetFlags  = log.SetFlags
	SetOutput = log.SetOutput
	SetPrefix = log.SetPrefix
	Writer    = log.Writer
)

var logger *Logger

func New(out io.Writer, prefix string, flag int) *Logger {
	log.SetOutput(out)
	logger = (*Logger)(log.New(out, prefix, flag))
	return logger
}

func NewRemoteSyslog(netType, addr, tag string) *Logger {
	fd, err := syslog.Dial(netType, addr, syslog.LOG_INFO, tag)
	if err != nil {
		return New(os.Stdout, "", log.LstdFlags)
	}
	return New(fd, "", log.LstdFlags)
}

func Debug(message string) {
	if logger != nil {
		logger.Debug(message)
	} else {
		log.Printf("[Debug]%s", message)
	}
}
func Info(message string) {
	if logger != nil {
		logger.Info(message)
	} else {
		log.Printf("[Info]%s", message)
	}
}
func Warn(message string) {
	if logger != nil {
		logger.Warn(message)
	} else {
		log.Printf("[Warn]%s", message)
	}
}
func Error(message string) {
	if logger != nil {
		logger.Error(message)
	} else {
		log.Printf("[Error]%s", message)
	}
}

func (l *Logger) Debug(message string) {
	l.format("Debug", message)
}
func (l *Logger) Info(message string) {
	l.format("Info", message)
}
func (l *Logger) Warn(message string) {
	l.format("Warn", message)
}
func (l *Logger) Error(message string) {
	l.format("Error", message)
}

func (l *Logger) format(level string, message string) {
	log.Println("[" + level + "] " + message)
}

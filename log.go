package log

import (
	"fmt"
	"io"
	"log"
	"log/syslog"
	"time"

	"github.com/panjf2000/ants/v2"
)

type Logger log.Logger

const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

var (
	Fatal                                      = log.Fatal
	Fatalf                                     = log.Fatalf
	Fatalln                                    = log.Fatalln
	Flags                                      = log.Flags
	Output                                     = log.Output
	Panic                                      = log.Panic
	Panicf                                     = log.Panicf
	Panicln                                    = log.Panicln
	Prefix                                     = log.Prefix
	Print                                      = log.Print
	Printf                                     = log.Printf
	Println                                    = log.Println
	SetFlags                                   = log.SetFlags
	SetOutput                                  = log.SetOutput
	SetPrefix                                  = log.SetPrefix
	Writer                                     = log.Writer
	DefaultPoolSize                            = 102400
	LevelString     map[syslog.Priority]string = map[syslog.Priority]string{
		syslog.LOG_ALERT:   "Alert",
		syslog.LOG_DEBUG:   "Debug",
		syslog.LOG_NOTICE:  "Notice",
		syslog.LOG_INFO:    "Info",
		syslog.LOG_WARNING: "Warn",
		syslog.LOG_ERR:     "Error",
		syslog.LOG_CRIT:    "Fatal",
		syslog.LOG_EMERG:   "Emerg",
	}
)

type logMessage struct {
	Message string
	Level   syslog.Priority
}

var logger *Logger

var pool *ants.PoolWithFunc

func New(out io.Writer, prefix string, flag int) *Logger {
	log.SetOutput(out)
	log.Default().SetFlags(log.LstdFlags | log.Llongfile | log.Lmicroseconds)
	logger = (*Logger)(log.New(out, prefix, flag))
	return logger
}

func NewRemoteSyslog(netType, addr, tag string, asyncsize int) (err error) {
	fd, err := syslog.Dial(netType, addr, syslog.LOG_DEBUG|syslog.LOG_KERN|syslog.LOG_WARNING|syslog.LOG_CRIT, tag)
	if err != nil {
		return
	}
	if asyncsize == 0 {
		asyncsize = DefaultPoolSize
	}
	pool, err = ants.NewPoolWithFunc(asyncsize, func(i interface{}) {
		payload, ok := i.(*logMessage)
		if !ok {
			return
		}
		switch payload.Level {
		case syslog.LOG_DEBUG:
			fd.Debug(payload.Message)
		case syslog.LOG_ALERT:
			if err := fd.Warning(payload.Message); err != nil {
				fmt.Println("syslog err", err)
			}
		case syslog.LOG_INFO:
			fd.Info(payload.Message)
		case syslog.LOG_WARNING:
			if err := fd.Warning(payload.Message); err != nil {
				fmt.Println("syslog err", err)
			}
		case syslog.LOG_ERR:
			fd.Err(payload.Message)
		case syslog.LOG_CRIT:
			fd.Crit(payload.Message)
		case syslog.LOG_EMERG:
			fd.Emerg(payload.Message)
		}
	}, ants.WithNonblocking(true), ants.WithExpiryDuration(5*time.Second))
	return
}

func ReleasePool() {
	if pool != nil {
		pool.Release()
		pool = nil
	}
}

func Debug(message string) {
	format(syslog.LOG_DEBUG, message)
}
func Info(message string) {
	format(syslog.LOG_INFO, message)
}
func Warn(message string) {
	format(syslog.LOG_WARNING, message)
}
func Error(message string) {
	format(syslog.LOG_ERR, message)
}

func (l *Logger) Debug(message string) {
	format(syslog.LOG_DEBUG, message)
}
func (l *Logger) Info(message string) {
	format(syslog.LOG_INFO, message)
}
func (l *Logger) Warn(message string) {
	format(syslog.LOG_WARNING, message)
}
func (l *Logger) Error(message string) {
}

func format(level syslog.Priority, message string) {
	if pool != nil {
		if err := pool.Invoke(&logMessage{
			Message: message,
			Level:   level,
		}); err != nil {
			fmt.Println(err)
		}
	} else {
		log.Println("[" + LevelString[level] + "] " + message)
	}
	//log.Println("[" + level + "] " + message)

	// fmt.Fprintf(*writer, "%s ["+level+"] %s", time.Now().Local().Format("2006-01-02T15:04:05.999Z"), message)
}

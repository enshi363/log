package log

import (
	"io"
	"log"
	"log/syslog"
	"os"
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
	log.Default().SetFlags(log.LstdFlags | log.Llongfile | log.Lmicroseconds)
	logger = (*Logger)(log.New(out, prefix, flag))
	return logger
}

func NewRemoteSyslog(netType, addr, tag string) *Logger {
	fd, err := syslog.Dial(netType, addr, syslog.LOG_DEBUG|syslog.LOG_KERN|syslog.LOG_WARNING|syslog.LOG_CRIT, tag)
	if err != nil {
		return New(os.Stdout, "", log.LstdFlags|log.Llongfile|log.Lmicroseconds)
	}
	return New(fd, "", log.LstdFlags|log.Llongfile|log.Lmicroseconds)
}

func Debug(message string) {
	if logger != nil {
		logger.Debug(message)
	} else {
		// fmt.Fprintf(*writer, "%s [Debug] %s", time.Now().Local().Format("2006-01-02T15:04:05.999Z"), message)
		log.Printf("[Debug] %s", message)
	}
}
func Info(message string) {
	if logger != nil {
		logger.Info(message)
	} else {
		// fmt.Fprintf(*writer, "%s [Info] %s", time.Now().Local().Format("2006-01-02T15:04:05.999Z"), message)
		log.Printf("[Info] %s", message)
	}
}
func Warn(message string) {
	if logger != nil {
		logger.Warn(message)
	} else {
		// fmt.Fprintf(*writer, "%s [Warn] %s", time.Now().Local().Format("2006-01-02T15:04:05.999Z"), message)
		log.Printf("[Warn] %s", message)
	}
}
func Error(message string) {
	if logger != nil {
		logger.Error(message)
	} else {
		// fmt.Fprintf(*writer, "%s [Error] %s", time.Now().Local().Format("2006-01-02T15:04:05.999Z"), message)
		log.Printf("[Error] %s", message)
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
	// fmt.Fprintf(*writer, "%s ["+level+"] %s", time.Now().Local().Format("2006-01-02T15:04:05.999Z"), message)
}

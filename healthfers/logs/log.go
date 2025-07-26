package logs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type AgreeGateLoager struct {
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
}

var (
	WarningLog *log.Logger
	InfoLog    *log.Logger
	ErrorLog   *log.Logger
)

func init() {
	file, err := os.OpenFile("myLOG.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLog = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLog = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLog = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func NewAgreeGateLogger() *AgreeGateLoager {
	logger := &AgreeGateLoager{}
	logger.InfoLogger = InfoLog
	logger.WarnLogger = WarningLog
	logger.ErrorLogger = ErrorLog
	return logger
}

// getCallerInfo returns the full file path, line number, and function name of the caller
func getCallerInfo(skip int) (string, int, string) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0, "unknown"
	}
	// Get full file path
	file = filepath.Clean(file)
	// Get function name
	fn := runtime.FuncForPC(pc).Name()
	// Simplify function name (remove package path)
	shortFn := fn[strings.LastIndex(fn, ".")+1:]
	return file, line, shortFn
}

// Info logs with full file path, line number, and function name
func (l *AgreeGateLoager) Info(v ...interface{}) {
	file, line, fn := getCallerInfo(2) // Skip 2 frames: getCallerInfo and Info
	msg := fmt.Sprintf("[%s:%d %s] %v", file, line, fn, fmt.Sprint(v...))
	l.InfoLogger.Println(msg)
}

// Warn logs with full file path, line number, and function name
func (l *AgreeGateLoager) Warn(v ...interface{}) {
	file, line, fn := getCallerInfo(2) // Skip 2 frames: getCallerInfo and Warn
	msg := fmt.Sprintf("[%s:%d %s] %v", file, line, fn, fmt.Sprint(v...))
	l.WarnLogger.Println(msg)
}

// Error logs with full file path, line number, and function name
func (l *AgreeGateLoager) Error(v ...interface{}) {
	file, line, fn := getCallerInfo(2) // Skip 2 frames: getCallerInfo and Error
	msg := fmt.Sprintf("[%s:%d %s] %v", file, line, fn, fmt.Sprint(v...))
	l.ErrorLogger.Println(msg)
}

// ErrorWithStack logs with full file path, line number, function name, and stack trace
func (l *AgreeGateLoager) ErrorWithStack(v ...interface{}) {
	file, line, fn := getCallerInfo(2)
	// Capture stack trace
	buf := make([]byte, 1<<16) // 64KB buffer
	stackSize := runtime.Stack(buf, false)
	stack := string(buf[:stackSize])
	msg := fmt.Sprintf("[%s:%d %s] %v\nStack Trace:\n%s", file, line, fn, fmt.Sprint(v...), stack)
	l.ErrorLogger.Println(msg)
}

// func (l *AgreeGateLoager) Info(v ...interface{}) {
// 	l.InfoLogger.Println(v...)
// }

// func (l *AgreeGateLoager) Warn(v ...interface{}) {
// 	l.WarnLogger.Println(v...)
// }

// func (l *AgreeGateLoager) Error(v ...interface{}) {
// 	l.ErrorLogger.Println(v...)
// }

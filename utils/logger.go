package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// 日志级别
const (
	DEBUG int = iota
	INFO
	WARN
	ERROR
	FATAL
)

var (
	// 默认日志级别
	logLevel = INFO
	// 日志前缀
	logPrefix = []string{
		"[DEBUG] ",
		"[INFO] ",
		"[WARN] ",
		"[ERROR] ",
		"[FATAL] ",
	}
	// 是否显示调用信息
	showCaller = true
	// 日志记录器
	infoLogger  = log.New(os.Stdout, "", log.LstdFlags)
	errorLogger = log.New(os.Stderr, "", log.LstdFlags)
)

// InitLogger 初始化日志系统
func InitLogger(level int, out io.Writer, errOut io.Writer, showCallerInfo bool) {
	logLevel = level
	if out != nil {
		infoLogger = log.New(out, "", log.LstdFlags)
	}
	if errOut != nil {
		errorLogger = log.New(errOut, "", log.LstdFlags)
	}
	showCaller = showCallerInfo
}

// getCallerInfo 获取调用者信息
func getCallerInfo() string {
	if !showCaller {
		return ""
	}

	_, file, line, ok := runtime.Caller(3) // 跳过getCallerInfo和日志函数本身
	if !ok {
		return ""
	}

	// 只使用文件名，不包含路径
	filename := filepath.Base(file)
	return fmt.Sprintf(" (%s:%d)", filename, line)
}

// Debug 打印调试日志
func Debug(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		callerInfo := getCallerInfo()
		infoLogger.Printf(logPrefix[DEBUG]+format+callerInfo, v...)
	}
}

// Info 打印信息日志
func Info(format string, v ...interface{}) {
	if logLevel <= INFO {
		callerInfo := getCallerInfo()
		infoLogger.Printf(logPrefix[INFO]+format+callerInfo, v...)
	}
}

// Warn 打印警告日志
func Warn(format string, v ...interface{}) {
	if logLevel <= WARN {
		callerInfo := getCallerInfo()
		infoLogger.Printf(logPrefix[WARN]+format+callerInfo, v...)
	}
}

// Error 打印错误日志
func Error(format string, v ...interface{}) {
	if logLevel <= ERROR {
		callerInfo := getCallerInfo()
		errorLogger.Printf(logPrefix[ERROR]+format+callerInfo, v...)
	}
}

// Fatal 打印致命错误日志并退出
func Fatal(format string, v ...interface{}) {
	if logLevel <= FATAL {
		callerInfo := getCallerInfo()
		errorLogger.Fatalf(logPrefix[FATAL]+format+callerInfo, v...)
	}
}

// InitLogFile 初始化日志文件
func InitLogFile(logDir string) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// 创建日志文件
	currentTime := time.Now().Format("2006-01-02")
	logFileName := filepath.Join(logDir, fmt.Sprintf("app_%s.log", currentTime))
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// 创建错误日志文件
	errLogFileName := filepath.Join(logDir, fmt.Sprintf("error_%s.log", currentTime))
	errLogFile, err := os.OpenFile(errLogFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}

	// 设置日志输出到文件和标准输出
	multiOut := io.MultiWriter(logFile, os.Stdout)
	multiErrOut := io.MultiWriter(errLogFile, os.Stderr)

	// 初始化日志系统
	InitLogger(logLevel, multiOut, multiErrOut, showCaller)
}

package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
)

// Logger 日志记录器结构体
type Logger struct {
	config *LogConfig
}

// LogConfig 日志配置
type LogConfig struct {
	LogDir     string
	LogLevel   zapcore.Level
	ShowCaller bool
}

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
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	// 初始化日志记录器
	infoLogger = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

// InitLogger 初始化日志系统.
func InitLogger(level int, stdout, stderr io.Writer, showCallerInfo bool) {
	// 设置日志级别.
	logLevel = level

	// 设置输出写入器.
	if stdout == nil {
		stdout = os.Stdout
	}

	infoLogger = log.New(stdout, "", log.LstdFlags)

	// 设置错误写入器.
	if stderr == nil {
		stderr = os.Stderr
	}

	errorLogger = log.New(stderr, "", log.LstdFlags)

	// 设置是否显示调用者信息.
	showCaller = showCallerInfo
}

// getCallerInfo 获取调用者信息.
func getCallerInfo() string {
	_, filename, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	return fmt.Sprintf(" (%s:%d)", filename, line)
}

// Debugf 打印调试日志
func Debugf(format string, v ...interface{}) {
	if logLevel <= DEBUG {
		callerInfo := getCallerInfo()
		infoLogger.Printf(logPrefix[DEBUG]+format+callerInfo, v...)
	}
}

// Infof 打印信息日志
func Infof(format string, v ...interface{}) {
	if logLevel <= INFO {
		callerInfo := getCallerInfo()
		infoLogger.Printf(logPrefix[INFO]+format+callerInfo, v...)
	}
}

// Warnf 打印警告日志
func Warnf(format string, v ...interface{}) {
	if logLevel <= WARN {
		callerInfo := getCallerInfo()
		infoLogger.Printf(logPrefix[WARN]+format+callerInfo, v...)
	}
}

// Errorf 打印错误日志
func Errorf(format string, v ...interface{}) {
	if logLevel <= ERROR {
		callerInfo := getCallerInfo()
		errorLogger.Printf(logPrefix[ERROR]+format+callerInfo, v...)
	}
}

// Fatalf 打印致命错误日志并退出
func Fatalf(format string, v ...interface{}) {
	if logLevel <= FATAL {
		callerInfo := getCallerInfo()
		errorLogger.Fatalf(logPrefix[FATAL]+format+callerInfo, v...)
	}
}

// InitLogFile 初始化日志文件
func InitLogFile(logDir string) {
	// 确保日志目录存在
	if err := os.MkdirAll(logDir, 0o750); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// 创建日志文件
	currentTime := time.Now().Format("2006-01-02")
	logFileName := filepath.Join(logDir, fmt.Sprintf("app_%s.log", currentTime))
	logFileName = cleanPath(logFileName)

	if !isPathSafe(logFileName, logDir) {
		log.Fatalf("Invalid log file path")
	}

	// #nosec G304 -- 路径已经过清理和验证
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// 创建错误日志文件
	errLogFileName := filepath.Join(logDir, fmt.Sprintf("error_%s.log", currentTime))
	errLogFileName = cleanPath(errLogFileName)

	if !isPathSafe(errLogFileName, logDir) {
		log.Fatalf("Invalid error log file path")
	}

	// #nosec G304 -- 路径已经过清理和验证
	errLogFile, err := os.OpenFile(errLogFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatalf("Failed to open error log file: %v", err)
	}

	// 设置日志输出到文件和标准输出
	multiOut := io.MultiWriter(logFile, os.Stdout)
	multiErrOut := io.MultiWriter(errLogFile, os.Stderr)

	// 初始化日志系统
	InitLogger(logLevel, multiOut, multiErrOut, showCaller)
}

// LogInfo 记录信息日志
func LogInfo(msg string, fields ...interface{}) {
	if len(fields)%2 != 0 {
		errorLogger.Printf("Invalid number of fields for LogInfo: %d", len(fields))
		return
	}

	// 构建日志消息
	logMsg := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), msg)
	fieldCount := len(fields)

	for i := 0; i < fieldCount; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			errorLogger.Printf("Invalid key type in LogInfo: %T", fields[i])
			continue
		}
		value := fields[i+1]
		logMsg = fmt.Sprintf("%s %s=%v", logMsg, key, value)
	}

	infoLogger.Println(logMsg)
}

// LogError 记录错误日志
func LogError(msg string, fields ...interface{}) {
	if len(fields)%2 != 0 {
		errorLogger.Printf("Invalid number of fields for LogError: %d", len(fields))
		return
	}

	// 构建日志消息
	logMsg := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05"), msg)
	fieldCount := len(fields)

	for i := 0; i < fieldCount; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			errorLogger.Printf("Invalid key type in LogError: %T", fields[i])
			continue
		}
		value := fields[i+1]
		logMsg = fmt.Sprintf("%s %s=%v", logMsg, key, value)
	}

	errorLogger.Println(logMsg)
}

// 添加路径清理函数
func cleanPath(path string) string {
	// 移除所有 .. 和 . 路径组件
	return filepath.Clean(path)
}

// 验证路径是否在允许的目录内
func isPathSafe(path, baseDir string) bool {
	cleanBase := filepath.Clean(baseDir)
	cleanPath := filepath.Clean(path)
	return strings.HasPrefix(cleanPath, cleanBase)
}

// NewLogger 创建新的日志记录器
func NewLogger(config *LogConfig) *Logger {
	return &Logger{
		config: config,
	}
}

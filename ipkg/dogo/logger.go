package dogo

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// LogLevel 日志级别
type LogLevel int

const (
	// LevelSilent 静默模式 - 只输出关键信息
	LevelSilent LogLevel = iota
	// LevelNormal 普通模式 - 输出所有信息
	LevelNormal
	// LevelVerbose 详细模式 - 输出调试信息
	LevelVerbose
)

// Logger 日志器
type Logger struct {
	mu     sync.Mutex
	level  LogLevel
	output io.Writer
}

// 全局日志实例
var (
	defaultLogger = NewLogger(LevelNormal, os.Stdout)
)

// NewLogger 创建新的日志器
func NewLogger(level LogLevel, output io.Writer) *Logger {
	return &Logger{
		level:  level,
		output: output,
	}
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetOutput 设置输出目标
func (l *Logger) SetOutput(output io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.output = output
}

// GetLevel 获取当前日志级别
func (l *Logger) GetLevel() LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level
}

// IsQuiet 是否静默模式
func (l *Logger) IsQuiet() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.level == LevelSilent
}

// Print 打印（受日志级别控制）
func (l *Logger) Print(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelNormal {
		fmt.Fprint(l.output, a...)
	}
}

// Println 打印并换行（受日志级别控制）
func (l *Logger) Println(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelNormal {
		fmt.Fprintln(l.output, a...)
	}
}

// Printf 格式化打印（受日志级别控制）
func (l *Logger) Printf(format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelNormal {
		fmt.Fprintf(l.output, format, a...)
	}
}

// Always 总是打印（忽略日志级别）
func (l *Logger) Always(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprint(l.output, a...)
}

// Alwaysln 总是打印并换行（忽略日志级别）
func (l *Logger) Alwaysln(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintln(l.output, a...)
}

// Alwaysf 总是格式化打印（忽略日志级别）
func (l *Logger) Alwaysf(format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	fmt.Fprintf(l.output, format, a...)
}

// Verbose 详细日志（仅在详细模式下输出）
func (l *Logger) Verbose(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelVerbose {
		fmt.Fprint(l.output, a...)
	}
}

// Verboseln 详细日志并换行（仅在详细模式下输出）
func (l *Logger) Verboseln(a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelVerbose {
		fmt.Fprintln(l.output, a...)
	}
}

// Verbosef 格式化详细日志（仅在详细模式下输出）
func (l *Logger) Verbosef(format string, a ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.level >= LevelVerbose {
		fmt.Fprintf(l.output, format, a...)
	}
}

// ============ 全局函数（使用默认日志器） ============

// SetQuietMode 设置静默模式
func SetQuietMode(quiet bool) {
	if quiet {
		defaultLogger.SetLevel(LevelSilent)
	} else {
		defaultLogger.SetLevel(LevelNormal)
	}
}

// SetLogLevel 设置日志级别
func SetLogLevel(level LogLevel) {
	defaultLogger.SetLevel(level)
}

// IsQuiet 是否静默模式
func IsQuiet() bool {
	return defaultLogger.IsQuiet()
}

// LogPrint 打印（受日志级别控制）
func LogPrint(a ...interface{}) {
	defaultLogger.Print(a...)
}

// LogPrintln 打印并换行（受日志级别控制）
func LogPrintln(a ...interface{}) {
	defaultLogger.Println(a...)
}

// LogPrintf 格式化打印（受日志级别控制）
func LogPrintf(format string, a ...interface{}) {
	defaultLogger.Printf(format, a...)
}

// LogAlways 总是打印（忽略日志级别）
func LogAlways(a ...interface{}) {
	defaultLogger.Always(a...)
}

// LogAlwaysln 总是打印并换行（忽略日志级别）
func LogAlwaysln(a ...interface{}) {
	defaultLogger.Alwaysln(a...)
}

// LogAlwaysf 总是格式化打印（忽略日志级别）
func LogAlwaysf(format string, a ...interface{}) {
	defaultLogger.Alwaysf(format, a...)
}

// LogVerbose 详细日志（仅在详细模式下输出）
func LogVerbose(a ...interface{}) {
	defaultLogger.Verbose(a...)
}

// LogVerboseln 详细日志并换行（仅在详细模式下输出）
func LogVerboseln(a ...interface{}) {
	defaultLogger.Verboseln(a...)
}

// LogVerbosef 格式化详细日志（仅在详细模式下输出）
func LogVerbosef(format string, a ...interface{}) {
	defaultLogger.Verbosef(format, a...)
}
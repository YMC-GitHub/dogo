package dogo

import (
	"time"
)

// Version 版本信息
const Version = "1.0.0"
const AppName = "dogo"

// Config 配置结构
type Config struct {
	Mirrors      []string `json:"mirrors"`
	Timeout      int      `json:"timeout"`       // 超时时间(秒)
	Concurrency  int      `json:"concurrency"`    // 并发数
	PrintResult  bool     `json:"print_result"`   // 打印测试结果
	PrintAdvice  bool     `json:"print_advice"`   // 打印优化建议
	PrintFails   bool     `json:"print_fails"`    // 打印不可用地址
	ShowProgress bool     `json:"show_progress"`  // 显示进度条
	Quiet        bool     `json:"quiet"`          // 静默模式
	ConfigFile   string   `json:"config_file"`    // 配置文件路径
	Format       string   `json:"format"`         // 输出格式: stro, strm
	Update       bool     `json:"update"`         // 更新输入文件
}

// TestResult 测试结果
type TestResult struct {
	Mirror    string        `json:"mirror"`
	Status    string        `json:"status"` // "success", "failed"
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error,omitempty"`
	CheckedAt time.Time     `json:"checked_at"`
}

// TestSummary 测试摘要
type TestSummary struct {
	Total        int           `json:"total"`
	Success      int           `json:"success"`
	Failed       int           `json:"failed"`
	AvgLatency   time.Duration `json:"avg_latency"`
	Fastest      []TestResult  `json:"fastest"`
	Available    []TestResult  `json:"available"`
	Unavailable  []TestResult  `json:"unavailable"`
}

// MirrorStats 镜像统计
type MirrorStats struct {
	URL       string        `json:"url"`
	Available bool          `json:"available"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error,omitempty"`
}
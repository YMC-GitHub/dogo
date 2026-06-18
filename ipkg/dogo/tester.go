package dogo

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Tester 测试器
type Tester struct {
	config   Config
	results  []TestResult
	mu       sync.Mutex
	progress int
	total    int
	stopChan chan struct{}
	logger   *Logger
}

// NewTester 创建测试器
func NewTester(config Config) *Tester {
	return &Tester{
		config:   config,
		results:  make([]TestResult, len(config.Mirrors)),
		total:    len(config.Mirrors),
		stopChan: make(chan struct{}),
		logger:   defaultLogger,
	}
}

// NewTesterWithLogger 创建带自定义日志器的测试器
func NewTesterWithLogger(config Config, logger *Logger) *Tester {
	tester := NewTester(config)
	tester.logger = logger
	return tester
}

// TestAll 测试所有镜像
func (t *Tester) TestAll() []TestResult {
	if t.config.ShowProgress && !t.logger.IsQuiet() {
		go t.showProgress()
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, t.config.Concurrency)

	for i, mirror := range t.config.Mirrors {
		wg.Add(1)
		go func(index int, m string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			result := t.TestSingle(m)
			t.mu.Lock()
			t.results[index] = result
			t.progress++
			t.mu.Unlock()
		}(i, mirror)
	}

	wg.Wait()
	
	if t.config.ShowProgress && !t.logger.IsQuiet() {
		fmt.Println()
	}
	
	return t.results
}

// showProgress 显示进度条
func (t *Tester) showProgress() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-t.stopChan:
			return
		case <-ticker.C:
			t.mu.Lock()
			progress := t.progress
			total := t.total
			t.mu.Unlock()
			
			if total > 0 {
				percent := float64(progress) / float64(total) * 100
				barWidth := 50
				filled := int(percent / 2)
				if filled > barWidth {
					filled = barWidth
				}
				
				bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)
				fmt.Printf("\r进度: [%s] %.1f%% (%d/%d)", bar, percent, progress, total)
				
				if progress >= total {
					return
				}
			}
		}
	}
}

// TestSingle 测试单个镜像
func (t *Tester) TestSingle(mirror string) TestResult {
	start := time.Now()
	cleanURL := strings.TrimSuffix(mirror, "/")
	
	client := &http.Client{
		Timeout: time.Duration(t.config.Timeout) * time.Second,
	}

	endpoints := []string{
		"/v2/",
		"/v2/_ping",
	}

	for _, endpoint := range endpoints {
		url := cleanURL + endpoint
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			continue
		}
		
		req.Header.Set("User-Agent", fmt.Sprintf("%s/%s", AppName, Version))

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		io.CopyN(io.Discard, resp.Body, 1024)

		if resp.StatusCode == 200 || resp.StatusCode == 401 {
			return TestResult{
				Mirror:    mirror,
				Status:    "success",
				Latency:   time.Since(start),
				CheckedAt: time.Now(),
			}
		}
	}

	return TestResult{
		Mirror:    mirror,
		Status:    "failed",
		Error:     "所有端点测试失败",
		CheckedAt: time.Now(),
	}
}

// GetProgress 获取测试进度
func (t *Tester) GetProgress() (int, int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.progress, t.total
}

// Stop 停止测试
func (t *Tester) Stop() {
	close(t.stopChan)
}
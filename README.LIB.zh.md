# dogo - Go 类库文档

dogo 是一个用于测试 Docker 镜像加速器可用性的 Go 类库，提供了一套完整的 API 用于批量测试、筛选和更新镜像配置。

## 📦 安装

```bash
go get github.com/ymc-github/dogo/ipkg/dogo
```

## 🏗️ 架构设计

```
dogo/ipkg/dogo/
├── types.go      # 类型定义 (Config, TestResult, TestSummary)
├── config.go     # 配置管理 (加载、保存、默认配置)
├── logger.go     # 日志系统 (多级别日志、静默模式)
├── tester.go     # 测试核心 (并发测试、进度控制)
└── utils.go      # 工具函数 (统计、格式化、过滤)
```

## 📚 核心类型

### Config - 配置结构
```go
type Config struct {
    Mirrors      []string // 镜像列表
    Timeout      int      // 超时时间(秒)
    Concurrency  int      // 并发数
    PrintResult  bool     // 打印测试结果
    PrintAdvice  bool     // 打印优化建议
    PrintFails   bool     // 打印不可用地址
    ShowProgress bool     // 显示进度条
    Quiet        bool     // 静默模式
    ConfigFile   string   // 配置文件路径
    Format       string   // 输出格式: stro, strm
    Update       bool     // 更新输入文件
}
```

### TestResult - 测试结果
```go
type TestResult struct {
    Mirror    string        // 镜像URL
    Status    string        // "success" 或 "failed"
    Latency   time.Duration // 响应延迟
    Error     string        // 错误信息
    CheckedAt time.Time     // 检查时间
}
```

### TestSummary - 测试摘要
```go
type TestSummary struct {
    Total        int           // 总数量
    Success      int           // 成功数量
    Failed       int           // 失败数量
    AvgLatency   time.Duration // 平均延迟
    Fastest      []TestResult  // 最快的镜像
    Available    []TestResult  // 可用镜像
    Unavailable  []TestResult  // 不可用镜像
}
```

## 🔧 API 参考

### 配置管理

#### DefaultConfig - 获取默认配置
```go
func DefaultConfig() Config
```
返回包含默认镜像列表和参数的配置。

**示例：**
```go
config := dogo.DefaultConfig()
config.Mirrors = []string{"https://mirror.aliyuncs.com"}
```

#### LoadConfigFromFile - 从文件加载配置
```go
func LoadConfigFromFile(path string) (Config, error)
```
从 JSON 文件加载配置，自动识别 `registry-mirrors` 或 `mirrors` 字段。

**示例：**
```go
config, err := dogo.LoadConfigFromFile("/etc/docker/daemon.json")
if err != nil {
    log.Fatal(err)
}
```

#### UpdateMirrorsInFile - 更新配置文件
```go
func UpdateMirrorsInFile(path string, mirrors []string) error
```
更新配置文件中的镜像列表，只修改指定键，保留其他配置。

**示例：**
```go
mirrors := []string{"https://mirror.aliyuncs.com", "https://docker.m.daocloud.io"}
err := dogo.UpdateMirrorsInFile("/etc/docker/daemon.json", mirrors)
```

### 测试器

#### NewTester - 创建测试器
```go
func NewTester(config Config) *Tester
```
使用指定配置创建测试器实例。

**示例：**
```go
config := dogo.DefaultConfig()
tester := dogo.NewTester(config)
```

#### TestAll - 测试所有镜像
```go
func (t *Tester) TestAll() []TestResult
```
并发测试所有配置的镜像，返回测试结果列表。

**示例：**
```go
results := tester.TestAll()
```

#### TestSingle - 测试单个镜像
```go
func (t *Tester) TestSingle(mirror string) TestResult
```
测试单个镜像的可用性。

**示例：**
```go
result := tester.TestSingle("https://mirror.aliyuncs.com")
if result.Status == "success" {
    fmt.Printf("可用，延迟: %v\n", result.Latency)
}
```

#### GetProgress - 获取测试进度
```go
func (t *Tester) GetProgress() (int, int)
```
获取当前测试进度（已完成数，总数）。

**示例：**
```go
done, total := tester.GetProgress()
fmt.Printf("进度: %d/%d\n", done, total)
```

#### Stop - 停止测试
```go
func (t *Tester) Stop()
```
停止正在进行的测试。

### 日志系统

#### 日志级别
```go
const (
    LevelSilent  LogLevel = iota // 静默模式
    LevelNormal                  // 普通模式
    LevelVerbose                 // 详细模式
)
```

#### SetQuietMode - 设置静默模式
```go
func SetQuietMode(quiet bool)
```
全局设置静默模式。

**示例：**
```go
dogo.SetQuietMode(true) // 启用静默模式
```

#### SetLogLevel - 设置日志级别
```go
func SetLogLevel(level LogLevel)
```
设置全局日志级别。

**示例：**
```go
dogo.SetLogLevel(dogo.LevelVerbose) // 启用详细日志
```

#### IsQuiet - 检查是否静默模式
```go
func IsQuiet() bool
```
返回当前是否处于静默模式。

### 工具函数

#### GetSummary - 获取测试摘要
```go
func GetSummary(results []TestResult) TestSummary
```
从测试结果生成统计摘要。

**示例：**
```go
results := tester.TestAll()
summary := dogo.GetSummary(results)
fmt.Printf("可用: %d/%d\n", summary.Success, summary.Total)
```

#### GetFastestMirrors - 获取最快的镜像
```go
func GetFastestMirrors(results []TestResult, n int) []TestResult
```
返回响应速度最快的 N 个镜像。

**示例：**
```go
fastest := dogo.GetFastestMirrors(results, 3)
for i, r := range fastest {
    fmt.Printf("%d. %s (%v)\n", i+1, r.Mirror, r.Latency)
}
```

#### GetAvailableMirrors - 获取可用镜像
```go
func GetAvailableMirrors(results []TestResult) []TestResult
```
返回所有可用的镜像。

#### GetUnavailableMirrors - 获取不可用镜像
```go
func GetUnavailableMirrors(results []TestResult) []TestResult
```
返回所有不可用的镜像。

#### ExtractMirrorURLs - 提取URL列表
```go
func ExtractMirrorURLs(results []TestResult) []string
```
从测试结果中提取可用的镜像URL列表。

**示例：**
```go
urls := dogo.ExtractMirrorURLs(results)
fmt.Println(strings.Join(urls, ","))
```

#### ExtractFailedURLs - 提取失败URL列表
```go
func ExtractFailedURLs(results []TestResult) []string
```
从测试结果中提取失败的镜像URL列表。

#### FormatMirrors - 格式化镜像列表
```go
func FormatMirrors(mirrors []string, format string) string
```
将镜像列表格式化为指定格式（stro 或 strm）。

**示例：**
```go
mirrors := []string{"url1", "url2"}
fmt.Println(dogo.FormatMirrors(mirrors, "stro")) // url1,url2
fmt.Println(dogo.FormatMirrors(mirrors, "strm")) // url1\nurl2
```

#### FormatMirrorsWithError - 带错误信息的格式化
```go
func FormatMirrorsWithError(results []TestResult, format string) string
```
格式化镜像列表，包含错误信息（用于调试）。

#### ValidateFormat - 验证格式
```go
func ValidateFormat(format string) bool
```
验证输出格式是否有效（stro 或 strm）。

#### NormalizeURL - 标准化URL
```go
func NormalizeURL(url string) string
```
标准化URL（去除空格和尾部斜杠）。

## 💡 使用示例

### 1. 基础用法 - 测试镜像列表

```go
package main

import (
    "fmt"
    "log"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    // 创建配置
    config := dogo.DefaultConfig()
    config.Mirrors = []string{
        "https://mirror.aliyuncs.com",
        "https://docker.m.daocloud.io",
        "https://docker.nju.edu.cn",
    }
    config.Timeout = 10
    config.Concurrency = 20

    // 创建测试器
    tester := dogo.NewTester(config)
    
    // 运行测试
    results := tester.TestAll()
    
    // 获取摘要
    summary := dogo.GetSummary(results)
    
    fmt.Printf("总测试: %d, 可用: %d, 不可用: %d\n", 
        summary.Total, summary.Success, summary.Failed)
    fmt.Printf("平均延迟: %v\n", summary.AvgLatency)
    
    // 显示最快的镜像
    fastest := dogo.GetFastestMirrors(results, 3)
    for i, r := range fastest {
        fmt.Printf("%d. %s (%v)\n", i+1, r.Mirror, r.Latency)
    }
}
```

### 2. 从配置文件加载并测试

```go
package main

import (
    "fmt"
    "log"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    // 从配置文件加载
    config, err := dogo.LoadConfigFromFile("/etc/docker/daemon.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // 覆盖部分配置
    config.Timeout = 5
    config.Concurrency = 10
    
    // 执行测试
    tester := dogo.NewTester(config)
    results := tester.TestAll()
    summary := dogo.GetSummary(results)
    
    // 输出可用镜像列表
    urls := dogo.ExtractMirrorURLs(results)
    fmt.Println("可用镜像:", urls)
    
    // 输出不可用镜像列表
    failed := dogo.ExtractFailedURLs(results)
    fmt.Println("不可用镜像:", failed)
}
```

### 3. 使用自定义日志器

```go
package main

import (
    "os"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    // 创建自定义日志器（输出到文件）
    file, _ := os.Create("test.log")
    logger := dogo.NewLogger(dogo.LevelNormal, file)
    
    // 配置
    config := dogo.DefaultConfig()
    config.Mirrors = []string{"https://mirror.aliyuncs.com"}
    
    // 使用自定义日志器
    tester := dogo.NewTesterWithLogger(config, logger)
    results := tester.TestAll()
    _ = results
}
```

### 4. 静默模式 + 格式化输出

```go
package main

import (
    "fmt"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    // 启用静默模式
    dogo.SetQuietMode(true)
    
    config := dogo.DefaultConfig()
    config.Mirrors = []string{
        "https://mirror.aliyuncs.com",
        "https://docker.m.daocloud.io",
    }
    
    tester := dogo.NewTester(config)
    results := tester.TestAll()
    summary := dogo.GetSummary(results)
    
    // 只输出可用的URL（单行格式）
    if len(summary.Available) > 0 {
        urls := dogo.ExtractMirrorURLs(summary.Available)
        fmt.Println(dogo.FormatMirrors(urls, "stro"))
    }
}
```

### 5. 测试并更新配置文件

```go
package main

import (
    "log"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    configFile := "/etc/docker/daemon.json"
    
    // 加载配置
    config, err := dogo.LoadConfigFromFile(configFile)
    if err != nil {
        log.Fatal(err)
    }
    
    // 执行测试
    tester := dogo.NewTester(config)
    results := tester.TestAll()
    summary := dogo.GetSummary(results)
    
    // 提取可用镜像
    available := dogo.ExtractMirrorURLs(summary.Available)
    if len(available) == 0 {
        log.Fatal("没有可用的镜像")
    }
    
    // 更新配置文件（只更新镜像列表）
    if err := dogo.UpdateMirrorsInFile(configFile, available); err != nil {
        log.Fatal(err)
    }
    
    log.Printf("已更新配置文件，保留了 %d 个可用镜像\n", len(available))
}
```

### 6. 并发测试 + 进度监控

```go
package main

import (
    "fmt"
    "time"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    config := dogo.DefaultConfig()
    tester := dogo.NewTester(config)
    
    // 在另一个goroutine监控进度
    go func() {
        for {
            done, total := tester.GetProgress()
            fmt.Printf("\r进度: %d/%d", done, total)
            if done >= total {
                break
            }
            time.Sleep(100 * time.Millisecond)
        }
        fmt.Println()
    }()
    
    // 执行测试
    results := tester.TestAll()
    summary := dogo.GetSummary(results)
    
    fmt.Printf("完成! 可用: %d/%d\n", summary.Success, summary.Total)
}
```

## 🧪 测试示例

```go
package dogo_test

import (
    "testing"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func TestDefaultConfig(t *testing.T) {
    config := dogo.DefaultConfig()
    if len(config.Mirrors) == 0 {
        t.Error("默认配置没有镜像列表")
    }
}

func TestValidateFormat(t *testing.T) {
    if !dogo.ValidateFormat("stro") {
        t.Error("stro 应该是有效格式")
    }
    if !dogo.ValidateFormat("strm") {
        t.Error("strm 应该是有效格式")
    }
    if dogo.ValidateFormat("invalid") {
        t.Error("invalid 应该是无效格式")
    }
}

func TestNormalizeURL(t *testing.T) {
    url := dogo.NormalizeURL(" https://example.com/ ")
    if url != "https://example.com" {
        t.Errorf("期望: https://example.com, 得到: %s", url)
    }
}
```

## 📊 性能建议

1. **并发数控制**：建议保持在 10-30 之间，避免触发限流
2. **超时设置**：根据网络环境调整，建议 5-15 秒
3. **批量测试**：大量镜像时建议分批测试
4. **缓存结果**：测试结果可缓存，避免重复测试

## 🔗 相关链接

- [GitHub 仓库](https://github.com/ymc-github/dogo)
- [命令行工具文档](README.md)
- [中文文档](README.zh.md)

## 📝 更新日志

### v1.0.0
- ✅ 初始版本发布
- ✅ 完整的 API 接口
- ✅ 并发测试支持
- ✅ 多格式输出
- ✅ 日志系统
- ✅ 配置管理
- ✅ 进度控制
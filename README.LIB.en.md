# dogo - Go Library Documentation

dogo is a Go library for testing Docker mirror accelerator availability, providing a complete set of APIs for batch testing, filtering, and updating mirror configurations.

## 📦 Installation

```bash
go get github.com/ymc-github/dogo/ipkg/dogo
```

## 🏗️ Architecture Design

```
dogo/ipkg/dogo/
├── types.go      # Type definitions (Config, TestResult, TestSummary)
├── config.go     # Configuration management (load, save, defaults)
├── logger.go     # Logging system (multi-level, silent mode)
├── tester.go     # Test core (concurrent testing, progress control)
└── utils.go      # Utility functions (statistics, formatting, filtering)
```

## 📚 Core Types

### Config - Configuration Structure
```go
type Config struct {
    Mirrors      []string // Mirror list
    Timeout      int      // Timeout in seconds
    Concurrency  int      // Number of concurrent tests
    PrintResult  bool     // Print test results
    PrintAdvice  bool     // Print optimization advice
    PrintFails   bool     // Print unavailable addresses
    ShowProgress bool     // Show progress bar
    Quiet        bool     // Silent mode
    ConfigFile   string   // Configuration file path
    Format       string   // Output format: stro, strm
    Update       bool     // Update input file
}
```

### TestResult - Test Result
```go
type TestResult struct {
    Mirror    string        // Mirror URL
    Status    string        // "success" or "failed"
    Latency   time.Duration // Response latency
    Error     string        // Error message
    CheckedAt time.Time     // Check time
}
```

### TestSummary - Test Summary
```go
type TestSummary struct {
    Total        int           // Total count
    Success      int           // Success count
    Failed       int           // Failed count
    AvgLatency   time.Duration // Average latency
    Fastest      []TestResult  // Fastest mirrors
    Available    []TestResult  // Available mirrors
    Unavailable  []TestResult  // Unavailable mirrors
}
```

## 🔧 API Reference

### Configuration Management

#### DefaultConfig - Get Default Configuration
```go
func DefaultConfig() Config
```
Returns a configuration with default mirror list and parameters.

**Example:**
```go
config := dogo.DefaultConfig()
config.Mirrors = []string{"https://mirror.aliyuncs.com"}
```

#### LoadConfigFromFile - Load Configuration from File
```go
func LoadConfigFromFile(path string) (Config, error)
```
Loads configuration from a JSON file, automatically detecting `registry-mirrors` or `mirrors` fields.

**Example:**
```go
config, err := dogo.LoadConfigFromFile("/etc/docker/daemon.json")
if err != nil {
    log.Fatal(err)
}
```

#### UpdateMirrorsInFile - Update Configuration File
```go
func UpdateMirrorsInFile(path string, mirrors []string) error
```
Updates the mirror list in the configuration file, modifying only the specified key while preserving other settings.

**Example:**
```go
mirrors := []string{"https://mirror.aliyuncs.com", "https://docker.m.daocloud.io"}
err := dogo.UpdateMirrorsInFile("/etc/docker/daemon.json", mirrors)
```

### Tester

#### NewTester - Create Tester
```go
func NewTester(config Config) *Tester
```
Creates a tester instance with the specified configuration.

**Example:**
```go
config := dogo.DefaultConfig()
tester := dogo.NewTester(config)
```

#### TestAll - Test All Mirrors
```go
func (t *Tester) TestAll() []TestResult
```
Concurrently tests all configured mirrors and returns the test results.

**Example:**
```go
results := tester.TestAll()
```

#### TestSingle - Test Single Mirror
```go
func (t *Tester) TestSingle(mirror string) TestResult
```
Tests the availability of a single mirror.

**Example:**
```go
result := tester.TestSingle("https://mirror.aliyuncs.com")
if result.Status == "success" {
    fmt.Printf("Available, latency: %v\n", result.Latency)
}
```

#### GetProgress - Get Test Progress
```go
func (t *Tester) GetProgress() (int, int)
```
Returns the current test progress (completed count, total count).

**Example:**
```go
done, total := tester.GetProgress()
fmt.Printf("Progress: %d/%d\n", done, total)
```

#### Stop - Stop Testing
```go
func (t *Tester) Stop()
```
Stops the ongoing test.

### Logging System

#### Log Levels
```go
const (
    LevelSilent  LogLevel = iota // Silent mode
    LevelNormal                  // Normal mode
    LevelVerbose                 // Verbose mode
)
```

#### SetQuietMode - Set Quiet Mode
```go
func SetQuietMode(quiet bool)
```
Globally sets quiet mode.

**Example:**
```go
dogo.SetQuietMode(true) // Enable quiet mode
```

#### SetLogLevel - Set Log Level
```go
func SetLogLevel(level LogLevel)
```
Sets the global log level.

**Example:**
```go
dogo.SetLogLevel(dogo.LevelVerbose) // Enable verbose logging
```

#### IsQuiet - Check Quiet Mode
```go
func IsQuiet() bool
```
Returns whether quiet mode is currently enabled.

### Utility Functions

#### GetSummary - Get Test Summary
```go
func GetSummary(results []TestResult) TestSummary
```
Generates a statistical summary from test results.

**Example:**
```go
results := tester.TestAll()
summary := dogo.GetSummary(results)
fmt.Printf("Available: %d/%d\n", summary.Success, summary.Total)
```

#### GetFastestMirrors - Get Fastest Mirrors
```go
func GetFastestMirrors(results []TestResult, n int) []TestResult
```
Returns the N fastest mirrors by response speed.

**Example:**
```go
fastest := dogo.GetFastestMirrors(results, 3)
for i, r := range fastest {
    fmt.Printf("%d. %s (%v)\n", i+1, r.Mirror, r.Latency)
}
```

#### GetAvailableMirrors - Get Available Mirrors
```go
func GetAvailableMirrors(results []TestResult) []TestResult
```
Returns all available mirrors.

#### GetUnavailableMirrors - Get Unavailable Mirrors
```go
func GetUnavailableMirrors(results []TestResult) []TestResult
```
Returns all unavailable mirrors.

#### ExtractMirrorURLs - Extract URL List
```go
func ExtractMirrorURLs(results []TestResult) []string
```
Extracts available mirror URLs from test results.

**Example:**
```go
urls := dogo.ExtractMirrorURLs(results)
fmt.Println(strings.Join(urls, ","))
```

#### ExtractFailedURLs - Extract Failed URL List
```go
func ExtractFailedURLs(results []TestResult) []string
```
Extracts failed mirror URLs from test results.

#### FormatMirrors - Format Mirror List
```go
func FormatMirrors(mirrors []string, format string) string
```
Formats the mirror list in the specified format (stro or strm).

**Example:**
```go
mirrors := []string{"url1", "url2"}
fmt.Println(dogo.FormatMirrors(mirrors, "stro")) // url1,url2
fmt.Println(dogo.FormatMirrors(mirrors, "strm")) // url1\nurl2
```

#### FormatMirrorsWithError - Format with Error Information
```go
func FormatMirrorsWithError(results []TestResult, format string) string
```
Formats the mirror list with error information (for debugging).

#### ValidateFormat - Validate Format
```go
func ValidateFormat(format string) bool
```
Validates whether the output format is valid (stro or strm).

#### NormalizeURL - Normalize URL
```go
func NormalizeURL(url string) string
```
Normalizes a URL (removes spaces and trailing slashes).

## 💡 Usage Examples

### 1. Basic Usage - Test Mirror List

```go
package main

import (
    "fmt"
    "log"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    // Create configuration
    config := dogo.DefaultConfig()
    config.Mirrors = []string{
        "https://mirror.aliyuncs.com",
        "https://docker.m.daocloud.io",
        "https://docker.nju.edu.cn",
    }
    config.Timeout = 10
    config.Concurrency = 20

    // Create tester
    tester := dogo.NewTester(config)
    
    // Run tests
    results := tester.TestAll()
    
    // Get summary
    summary := dogo.GetSummary(results)
    
    fmt.Printf("Total: %d, Available: %d, Unavailable: %d\n", 
        summary.Total, summary.Success, summary.Failed)
    fmt.Printf("Average latency: %v\n", summary.AvgLatency)
    
    // Display fastest mirrors
    fastest := dogo.GetFastestMirrors(results, 3)
    for i, r := range fastest {
        fmt.Printf("%d. %s (%v)\n", i+1, r.Mirror, r.Latency)
    }
}
```

### 2. Load and Test from Configuration File

```go
package main

import (
    "fmt"
    "log"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    // Load from configuration file
    config, err := dogo.LoadConfigFromFile("/etc/docker/daemon.json")
    if err != nil {
        log.Fatal(err)
    }
    
    // Override some settings
    config.Timeout = 5
    config.Concurrency = 10
    
    // Run tests
    tester := dogo.NewTester(config)
    results := tester.TestAll()
    summary := dogo.GetSummary(results)
    
    // Output available mirror list
    urls := dogo.ExtractMirrorURLs(results)
    fmt.Println("Available mirrors:", urls)
    
    // Output unavailable mirror list
    failed := dogo.ExtractFailedURLs(results)
    fmt.Println("Unavailable mirrors:", failed)
}
```

### 3. Using Custom Logger

```go
package main

import (
    "os"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    // Create custom logger (output to file)
    file, _ := os.Create("test.log")
    logger := dogo.NewLogger(dogo.LevelNormal, file)
    
    // Configuration
    config := dogo.DefaultConfig()
    config.Mirrors = []string{"https://mirror.aliyuncs.com"}
    
    // Use custom logger
    tester := dogo.NewTesterWithLogger(config, logger)
    results := tester.TestAll()
    _ = results
}
```

### 4. Silent Mode + Formatted Output

```go
package main

import (
    "fmt"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    // Enable silent mode
    dogo.SetQuietMode(true)
    
    config := dogo.DefaultConfig()
    config.Mirrors = []string{
        "https://mirror.aliyuncs.com",
        "https://docker.m.daocloud.io",
    }
    
    tester := dogo.NewTester(config)
    results := tester.TestAll()
    summary := dogo.GetSummary(results)
    
    // Only output available URLs (single-line format)
    if len(summary.Available) > 0 {
        urls := dogo.ExtractMirrorURLs(summary.Available)
        fmt.Println(dogo.FormatMirrors(urls, "stro"))
    }
}
```

### 5. Test and Update Configuration File

```go
package main

import (
    "log"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
    configFile := "/etc/docker/daemon.json"
    
    // Load configuration
    config, err := dogo.LoadConfigFromFile(configFile)
    if err != nil {
        log.Fatal(err)
    }
    
    // Run tests
    tester := dogo.NewTester(config)
    results := tester.TestAll()
    summary := dogo.GetSummary(results)
    
    // Extract available mirrors
    available := dogo.ExtractMirrorURLs(summary.Available)
    if len(available) == 0 {
        log.Fatal("No available mirrors found")
    }
    
    // Update configuration file (only update mirror list)
    if err := dogo.UpdateMirrorsInFile(configFile, available); err != nil {
        log.Fatal(err)
    }
    
    log.Printf("Configuration updated, kept %d available mirrors\n", len(available))
}
```

### 6. Concurrent Testing with Progress Monitoring

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
    
    // Monitor progress in another goroutine
    go func() {
        for {
            done, total := tester.GetProgress()
            fmt.Printf("\rProgress: %d/%d", done, total)
            if done >= total {
                break
            }
            time.Sleep(100 * time.Millisecond)
        }
        fmt.Println()
    }()
    
    // Run tests
    results := tester.TestAll()
    summary := dogo.GetSummary(results)
    
    fmt.Printf("Done! Available: %d/%d\n", summary.Success, summary.Total)
}
```

## 🧪 Testing Examples

```go
package dogo_test

import (
    "testing"
    "github.com/ymc-github/dogo/ipkg/dogo"
)

func TestDefaultConfig(t *testing.T) {
    config := dogo.DefaultConfig()
    if len(config.Mirrors) == 0 {
        t.Error("Default configuration has no mirror list")
    }
}

func TestValidateFormat(t *testing.T) {
    if !dogo.ValidateFormat("stro") {
        t.Error("stro should be a valid format")
    }
    if !dogo.ValidateFormat("strm") {
        t.Error("strm should be a valid format")
    }
    if dogo.ValidateFormat("invalid") {
        t.Error("invalid should be an invalid format")
    }
}

func TestNormalizeURL(t *testing.T) {
    url := dogo.NormalizeURL(" https://example.com/ ")
    if url != "https://example.com" {
        t.Errorf("Expected: https://example.com, Got: %s", url)
    }
}
```

## 📊 Performance Recommendations

1. **Concurrency Control**: Keep between 10-30 to avoid rate limiting
2. **Timeout Settings**: Adjust based on network conditions, recommended 5-15 seconds
3. **Batch Testing**: Consider testing in batches for large mirror lists
4. **Result Caching**: Cache test results to avoid repeated testing

## 🔗 Related Links

- [GitHub Repository](https://github.com/ymc-github/dogo)
- [CLI Tool Documentation](README.md)
- [Chinese Documentation](README.zh.md)

## 📝 Changelog

### v1.0.0
- ✅ Initial release
- ✅ Complete API interface
- ✅ Concurrent testing support
- ✅ Multi-format output
- ✅ Logging system
- ✅ Configuration management
- ✅ Progress control
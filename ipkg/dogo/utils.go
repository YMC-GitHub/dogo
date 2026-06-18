package dogo

import (
	"sort"
	"strings"
	"time"
)

// GetSummary 获取测试摘要
func GetSummary(results []TestResult) TestSummary {
	var summary TestSummary
	summary.Total = len(results)

	var available []TestResult
	var unavailable []TestResult
	var totalLatency time.Duration

	for _, r := range results {
		if r.Status == "success" {
			summary.Success++
			totalLatency += r.Latency
			available = append(available, r)
		} else {
			summary.Failed++
			unavailable = append(unavailable, r)
		}
	}

	summary.Available = available
	summary.Unavailable = unavailable

	if summary.Success > 0 {
		summary.AvgLatency = totalLatency / time.Duration(summary.Success)
		summary.Fastest = GetFastestMirrors(results, 5)
	}

	return summary
}

// GetFastestMirrors 获取最快的N个镜像
func GetFastestMirrors(results []TestResult, n int) []TestResult {
	var success []TestResult
	for _, r := range results {
		if r.Status == "success" {
			success = append(success, r)
		}
	}

	sort.Slice(success, func(i, j int) bool {
		return success[i].Latency < success[j].Latency
	})

	if n > len(success) {
		n = len(success)
	}
	return success[:n]
}

// GetAvailableMirrors 获取可用镜像列表
func GetAvailableMirrors(results []TestResult) []TestResult {
	var available []TestResult
	for _, r := range results {
		if r.Status == "success" {
			available = append(available, r)
		}
	}
	return available
}

// GetUnavailableMirrors 获取不可用镜像列表
func GetUnavailableMirrors(results []TestResult) []TestResult {
	var unavailable []TestResult
	for _, r := range results {
		if r.Status == "failed" {
			unavailable = append(unavailable, r)
		}
	}
	return unavailable
}

// FormatMirrors 格式化镜像列表
func FormatMirrors(mirrors []string, format string) string {
	if len(mirrors) == 0 {
		return ""
	}

	switch format {
	case "stro":
		return strings.Join(mirrors, ",")
	case "strm":
		return strings.Join(mirrors, "\n")
	default:
		return strings.Join(mirrors, ",")
	}
}

// FormatMirrorsWithError 格式化镜像列表（包含错误信息）
func FormatMirrorsWithError(results []TestResult, format string) string {
	if len(results) == 0 {
		return ""
	}

	var builder strings.Builder
	
	switch format {
	case "stro":
		var items []string
		for _, r := range results {
			if r.Error != "" {
				items = append(items, r.Mirror+"("+r.Error+")")
			} else {
				items = append(items, r.Mirror)
			}
		}
		builder.WriteString(strings.Join(items, ","))
		
	case "strm":
		for _, r := range results {
			if r.Error != "" {
				builder.WriteString(r.Mirror + " " + r.Error + "\n")
			} else {
				builder.WriteString(r.Mirror + "\n")
			}
		}
		
	default:
		for _, r := range results {
			if r.Error != "" {
				builder.WriteString(r.Mirror + " (错误: " + r.Error + ")\n")
			} else {
				builder.WriteString(r.Mirror + "\n")
			}
		}
	}
	
	return builder.String()
}

func FormatMirrorsOnly(results []TestResult, format string) string {
    var urls []string
    for _, r := range results {
        urls = append(urls, r.Mirror)
    }
    return FormatMirrors(urls, format)
}

// ExtractMirrorURLs 从测试结果提取URL列表
func ExtractMirrorURLs(results []TestResult) []string {
	var urls []string
	for _, r := range results {
		if r.Status == "success" {
			urls = append(urls, r.Mirror)
		}
	}
	return urls
}

// ExtractFailedURLs 从测试结果提取失败的URL列表
func ExtractFailedURLs(results []TestResult) []string {
	var urls []string
	for _, r := range results {
		if r.Status == "failed" {
			urls = append(urls, r.Mirror)
		}
	}
	return urls
}

// FilterMirrors 过滤镜像
func FilterMirrors(mirrors []string, available bool, results []TestResult) []string {
	availableMap := make(map[string]bool)
	for _, r := range results {
		if r.Status == "success" {
			availableMap[r.Mirror] = true
		}
	}

	var filtered []string
	for _, m := range mirrors {
		if availableMap[m] == available {
			filtered = append(filtered, m)
		}
	}
	return filtered
}

// ValidateFormat 验证格式
func ValidateFormat(format string) bool {
	return format == "stro" || format == "strm"
}

// NormalizeURL 标准化URL
func NormalizeURL(url string) string {
	url = strings.TrimSpace(url)
	url = strings.TrimSuffix(url, "/")
	return url
}
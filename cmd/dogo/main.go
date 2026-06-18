package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ymc-github/dogo/ipkg/dogo"
)

func main() {
	help := flag.Bool("help", false, "显示帮助信息")
	version := flag.Bool("version", false, "显示版本信息")
	file := flag.String("file", "", "配置文件路径 (必需)")
	printResult := flag.Bool("print-result", true, "打印测试结果")
	printAdviceFlag := flag.Bool("print-advice", true, "打印优化建议")
	printFails := flag.Bool("print-fails", false, "打印不可用地址")
	showProgress := flag.Bool("progress", true, "显示进度条")
	quiet := flag.Bool("quiet", false, "静默模式（只输出关键信息）")
	quietShort := flag.Bool("q", false, "静默模式（简写）")
	timeout := flag.Int("timeout", 10, "超时时间(秒)")
	concurrency := flag.Int("concurrency", 20, "并发数")
	format := flag.String("format", "stro", "输出格式: stro(单行), strm(多行)")
	update := flag.Bool("update", false, "更新输入文件")

	flag.Parse()

	if *quiet || *quietShort {
		dogo.SetQuietMode(true)
	}

	if *help || len(os.Args) > 1 && os.Args[1] == "help" {
		showHelp()
		return
	}
	if *version || len(os.Args) > 1 && os.Args[1] == "version" {
		showVersion()
		return
	}

	if *file == "" {
		dogo.LogAlwaysln("❌ 错误: --file 参数是必需的")
		dogo.LogAlwaysln("使用 --help 查看帮助")
		os.Exit(1)
	}

	if _, err := os.Stat(*file); err != nil {
		dogo.LogAlwaysf("❌ 错误: 配置文件不存在: %s\n", *file)
		os.Exit(1)
	}

	config, err := dogo.LoadConfigFromFile(*file)
	if err != nil {
		dogo.LogAlwaysf("❌ 读取配置文件失败: %v\n", err)
		os.Exit(1)
	}

	config.ConfigFile = *file
	config.Timeout = *timeout
	config.Concurrency = *concurrency
	config.PrintResult = *printResult
	config.PrintAdvice = *printAdviceFlag
	config.PrintFails = *printFails
	config.ShowProgress = *showProgress
	config.Quiet = dogo.IsQuiet()
	config.Format = *format
	config.Update = *update

	if !dogo.ValidateFormat(config.Format) {
		if !dogo.IsQuiet() {
			dogo.LogPrintf("⚠️  不支持的格式: %s，使用默认格式 stro\n", config.Format)
		}
		config.Format = "stro"
	}

	if !dogo.IsQuiet() {
		dogo.LogPrintf("🚀 开始测试 %d 个镜像加速器...\n", len(config.Mirrors))
		dogo.LogPrintf("📁 配置文件: %s\n", config.ConfigFile)
		dogo.LogPrintf("⏱️  超时: %ds, 并发数: %d\n", config.Timeout, config.Concurrency)
		if config.ShowProgress {
			dogo.LogPrintln("📊 进度条已启用")
		}
		dogo.LogPrintln()
	}

	tester := dogo.NewTester(config)
	results := tester.TestAll()
	summary := dogo.GetSummary(results)

	if config.Update {
		updateConfigFile(config, summary)
	}

	if config.PrintResult && !dogo.IsQuiet() {
		printResults(summary)
	}

	if config.PrintAdvice && !dogo.IsQuiet() {
		printAdviceFunc(summary)
	}

	if len(summary.Available) > 0 {
		printAvailableMirrors(summary, config.Format)
	}

	if config.PrintFails {
		printFailedMirrors(summary, config.Format)
	}

	if dogo.IsQuiet() {
		fmt.Printf("\n可用: %d/%d, 平均延迟: %v\n", 
			summary.Success, 
			summary.Total, 
			summary.AvgLatency.Round(time.Millisecond))
	}
}

func printAvailableMirrors(summary dogo.TestSummary, format string) {
	if len(summary.Available) == 0 {
		return
	}
	
	if dogo.IsQuiet() {
		urls := dogo.ExtractMirrorURLs(summary.Available)
		fmt.Println(dogo.FormatMirrors(urls, format))
		return
	}
	
	fmt.Println("\n📋 可用镜像列表:")
	fmt.Println(strings.Repeat("-", 80))
	
	urls := dogo.ExtractMirrorURLs(summary.Available)
	fmt.Println(dogo.FormatMirrors(urls, format))
}

func printFailedMirrors(summary dogo.TestSummary, format string) {
	if len(summary.Unavailable) == 0 {
		if !dogo.IsQuiet() {
			fmt.Println("\n✅ 所有镜像都可用！")
		}
		return
	}
	
	// 对于 print-fails，只输出 URL 列表，不包含错误信息
	if dogo.IsQuiet() {
		// 静默模式下，只输出URL列表（方便脚本使用）
		urls := dogo.ExtractFailedURLs(summary.Unavailable)
		fmt.Println(dogo.FormatMirrors(urls, format))
		return
	}
	
	// 非静默模式下，显示标题和URL列表（不包含错误信息）
	fmt.Println("\n❌ 不可用镜像列表:")
	fmt.Println(strings.Repeat("-", 80))
	
	// 只输出URL，不包含错误信息
	urls := dogo.ExtractFailedURLs(summary.Unavailable)
	fmt.Println(dogo.FormatMirrors(urls, format))
	
	// 显示统计信息
	fmt.Printf("\n📊 不可用镜像: %d 个\n", len(summary.Unavailable))
}

func printResults(summary dogo.TestSummary) {
	fmt.Println("\n📊 测试结果:")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range summary.Available {
		fmt.Printf("✅ %-50s 延迟: %v\n", r.Mirror, r.Latency.Round(time.Millisecond))
	}
	for _, r := range summary.Unavailable {
		fmt.Printf("❌ %-50s %s\n", r.Mirror, r.Error)
	}

	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("📈 统计: %d/%d 可用, 平均延迟: %v\n", 
		summary.Success, summary.Total, summary.AvgLatency.Round(time.Millisecond))

	if len(summary.Fastest) > 0 {
		fmt.Println("\n⚡ 最快镜像 (Top 3):")
		for i, r := range summary.Fastest[:min(3, len(summary.Fastest))] {
			fmt.Printf("  %d. %s (延迟: %v)\n", i+1, r.Mirror, r.Latency.Round(time.Millisecond))
		}
	}
}

func printAdviceFunc(summary dogo.TestSummary) {
	if len(summary.Fastest) == 0 {
		fmt.Println("\n⚠️  没有可用的镜像加速器")
		return
	}

	fmt.Println("\n💡 优化建议:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Println("推荐使用的镜像加速器（按速度排序）:")
	for i, r := range summary.Fastest[:min(5, len(summary.Fastest))] {
		fmt.Printf("  %d. %s (延迟: %v)\n", i+1, r.Mirror, r.Latency.Round(time.Millisecond))
	}
}

func updateConfigFile(config dogo.Config, summary dogo.TestSummary) {
	available := dogo.ExtractMirrorURLs(summary.Available)
	if len(available) == 0 {
		if !dogo.IsQuiet() {
			dogo.LogPrintln("\n⚠️  没有可用的镜像，无法更新配置文件")
		} else {
			dogo.LogAlwaysln("⚠️  没有可用的镜像，无法更新")
		}
		return
	}
	
	if !dogo.IsQuiet() {
		dogo.LogPrintf("\n📝 更新配置文件: %s\n", config.ConfigFile)
		dogo.LogPrintf("找到 %d 个可用镜像\n", len(available))
	}
	
	if err := dogo.UpdateMirrorsInFile(config.ConfigFile, available); err != nil {
		if !dogo.IsQuiet() {
			dogo.LogPrintf("❌ 更新配置文件失败: %v\n", err)
		} else {
			dogo.LogAlwaysf("❌ 更新失败: %v\n", err)
		}
	} else {
		if !dogo.IsQuiet() {
			dogo.LogPrintf("✅ 配置文件已更新，保留了 %d 个可用镜像\n", len(available))
			dogo.LogPrintln("   (仅更新了镜像列表，其他配置保持不变)")
		} else {
			dogo.LogAlwaysf("✅ 已更新 %d 个可用镜像\n", len(available))
		}
	}
}

func showHelp() {
	fmt.Printf(`dogo - Docker镜像加速器测试工具

用法:
  dogo --file <path> [选项]

必需参数:
  --file <path>, -f       指定 daemon.json 配置文件路径 (必需)

选项:
  --help, -h              显示帮助信息
  --version, -v           显示版本信息
  --print-result          打印测试结果 (默认: true)
  --print-advice          打印优化建议 (默认: true)
  --print-fails           打印不可用地址
  --progress              显示进度条 (默认: true)
  --quiet, -q             静默模式（只输出关键信息）
  --timeout <seconds>     设置超时时间 (默认: 10秒)
  --concurrency <number>  设置并发数 (默认: 20)
  --format <format>       输出格式: stro(单行), strm(多行)
  --update                更新输入文件

示例:
  dogo --file /etc/docker/daemon.json
  dogo --file /etc/docker/daemon.json --progress=false
  dogo -f daemon.json -q --print-fails --format stro
  dogo --file config.json --update --quiet --progress=false
`)
}

func showVersion() {
	fmt.Printf("%s\n", dogo.Version)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
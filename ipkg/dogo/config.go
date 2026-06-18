package dogo

import (
	"encoding/json"
	"fmt"
	"os"
)

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Mirrors: []string{
			"https://docker.registry.cyou",
			"https://docker-cf.registry.cyou",
			"https://dockercf.jsdelivr.fyi",
			"https://docker.jsdelivr.fyi",
			"https://dockertest.jsdelivr.fyi",
			"https://mirror.aliyuncs.com",
			"https://dockerproxy.com",
			"https://mirror.baidubce.com",
			"https://docker.m.daocloud.io",
			"https://docker.nju.edu.cn",
			"https://docker.mirrors.sjtug.sjtu.edu.cn",
			"https://docker.mirrors.ustc.edu.cn",
			"https://mirror.iscas.ac.cn",
			"https://docker.rainbond.cc",
			"https://do.nark.eu.org",
			"https://dc.j8.work",
			"https://gst6rzl9.mirror.aliyuncs.com",
			"https://registry.docker-cn.com",
			"http://hub-mirror.c.163.com",
			"http://mirrors.ustc.edu.cn/",
			"https://mirrors.tuna.tsinghua.edu.cn/",
			"http://mirrors.sohu.com/",
		},
		Timeout:      10,
		Concurrency:  20,
		PrintResult:  true,
		PrintAdvice:  true,
		PrintFails:   false,
		ShowProgress: true,
		Quiet:        false,
		Format:       "stro",
		Update:       false,
	}
}

// LoadConfigFromFile 从文件加载配置
func LoadConfigFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("读取文件失败: %w", err)
	}

	var rawConfig map[string]interface{}
	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return Config{}, fmt.Errorf("解析JSON失败: %w", err)
	}

	config := DefaultConfig()
	
	mirrors := extractMirrors(rawConfig)
	if len(mirrors) > 0 {
		config.Mirrors = mirrors
	} else {
		return Config{}, fmt.Errorf("配置文件中没有找到 registry-mirrors 或 mirrors 字段")
	}

	extractOtherConfig(&config, rawConfig)
	
	return config, nil
}

// extractMirrors 提取镜像列表
func extractMirrors(rawConfig map[string]interface{}) []string {
	var mirrors []string
	
	if val, ok := rawConfig["registry-mirrors"]; ok {
		if list, ok := val.([]interface{}); ok {
			for _, item := range list {
				if str, ok := item.(string); ok {
					mirrors = append(mirrors, str)
				}
			}
		}
	}
	
	if len(mirrors) == 0 {
		if val, ok := rawConfig["mirrors"]; ok {
			if list, ok := val.([]interface{}); ok {
				for _, item := range list {
					if str, ok := item.(string); ok {
						mirrors = append(mirrors, str)
					}
				}
			}
		}
	}
	
	return mirrors
}

// extractOtherConfig 提取其他配置项
func extractOtherConfig(config *Config, rawConfig map[string]interface{}) {
	if val, ok := rawConfig["timeout"]; ok {
		if timeout, ok := val.(float64); ok {
			config.Timeout = int(timeout)
		}
	}
	if val, ok := rawConfig["concurrency"]; ok {
		if concurrency, ok := val.(float64); ok {
			config.Concurrency = int(concurrency)
		}
	}
	if val, ok := rawConfig["print_result"]; ok {
		if printResult, ok := val.(bool); ok {
			config.PrintResult = printResult
		}
	}
	if val, ok := rawConfig["print_advice"]; ok {
		if printAdvice, ok := val.(bool); ok {
			config.PrintAdvice = printAdvice
		}
	}
	if val, ok := rawConfig["print_fails"]; ok {
		if printFails, ok := val.(bool); ok {
			config.PrintFails = printFails
		}
	}
	if val, ok := rawConfig["show_progress"]; ok {
		if showProgress, ok := val.(bool); ok {
			config.ShowProgress = showProgress
		}
	}
	if val, ok := rawConfig["quiet"]; ok {
		if quiet, ok := val.(bool); ok {
			config.Quiet = quiet
		}
	}
	if val, ok := rawConfig["format"]; ok {
		if format, ok := val.(string); ok {
			config.Format = format
		}
	}
	if val, ok := rawConfig["update"]; ok {
		if update, ok := val.(bool); ok {
			config.Update = update
		}
	}
}

// UpdateMirrorsInFile 更新配置文件中的镜像列表（只更新指定键）
func UpdateMirrorsInFile(path string, mirrors []string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	var configMap map[string]interface{}
	if err := json.Unmarshal(data, &configMap); err != nil {
		return fmt.Errorf("解析JSON失败: %w", err)
	}

	keyName := detectMirrorKey(configMap)
	
	mirrorsInterface := make([]interface{}, len(mirrors))
	for i, m := range mirrors {
		mirrorsInterface[i] = m
	}

	configMap[keyName] = mirrorsInterface

	output, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %w", err)
	}

	if err := os.WriteFile(path, output, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	return nil
}

// detectMirrorKey 检测镜像列表的键名
func detectMirrorKey(configMap map[string]interface{}) string {
	if _, ok := configMap["registry-mirrors"]; ok {
		return "registry-mirrors"
	}
	if _, ok := configMap["mirrors"]; ok {
		return "mirrors"
	}
	return "mirrors"
}
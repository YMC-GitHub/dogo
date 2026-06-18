# dogo - Docker镜像加速器测试工具

自动化测试 Docker daemon.json 中的镜像加速器可用性，支持批量测试、智能筛选、自动更新配置。

## ✨ 核心功能

1. **批量镜像测试**：自动测试配置文件中所有镜像加速器的连通性和响应速度
2. **智能筛选排序**：按响应速度排序，自动识别可用/不可用镜像
3. **自动更新配置**：测试通过后自动更新配置文件，仅保留可用镜像
4. **多种输出格式**：支持单行逗号分隔（stro）和多行换行分隔（strm），方便脚本解析
5. **灵活显示控制**：可分别控制结果输出、优化建议、失败列表、进度条显示
6. **静默运行模式**：支持静默输出，适合CI/CD和自动化脚本
7. **安全预览**：支持 `--print-fails` 预览不可用镜像列表
8. **调试日志**：完善的日志系统，支持详细模式调试

## 📌 完整参数说明

### 必需参数
| 参数 | 短别名 | 说明 | 示例 |
|------|--------|------|------|
| --file | -f | 目标 JSON 文件路径（必需） | /etc/docker/daemon.json、daemon.json |

### 显示控制参数
| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| --print-result | 打印详细测试结果 | true | --print-result=false |
| --print-advice | 打印优化建议 | true | --print-advice=false |
| --print-fails | 打印不可用镜像列表 | false | --print-fails |
| --progress | 显示进度条 | true | --progress=false |
| --quiet, -q | 静默模式（只输出关键信息） | false | --quiet 或 -q |
| --format | 输出格式：stro(单行), strm(多行) | stro | --format strm |

### 测试控制参数
| 参数 | 说明 | 默认值 | 示例 |
|------|------|--------|------|
| --timeout | 单个镜像测试超时时间(秒) | 10 | --timeout 5 |
| --concurrency | 并发测试数量 | 20 | --concurrency 10 |
| --update | 测试完成后更新配置文件 | false | --update |

### 信息查看参数
| 参数 | 短别名 | 说明 |
|------|--------|------|
| --version | -v | 显示版本信息 |
| --help | -h | 显示完整帮助文档 |

## 🚀 最常用命令示例

### 1. 基础信息查看
```bash
# 查看版本
bin/dogo --version
bin/dogo -v

# 完整帮助文档
bin/dogo --help
bin/dogo -h
```

### 2. Docker daemon.json 测试（最典型场景）

#### 基础测试（显示所有信息）
```bash
sudo bin/dogo --file /etc/docker/daemon.json
```

#### 测试并自动更新配置文件（仅保留可用镜像）
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --update
```

#### 静默测试，只输出可用镜像列表（适合脚本）
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --quiet \
  --format stro \
  --print-result=false \
  --print-advice=false
```

#### 查看不可用镜像列表（不包含错误信息）
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --print-fails \
  --format stro
```

#### 静默查看不可用镜像（纯输出，适合脚本）
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --quiet \
  --print-fails \
  --format stro \
  --print-result=false \
  --print-advice=false
```

#### 测试完成后重载Docker生效
```bash
sudo systemctl daemon-reload
sudo systemctl restart docker
```

### 3. 自定义测试配置

#### 调整超时时间和并发数
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --timeout 5 \
  --concurrency 10
```

#### 隐藏进度条，只显示结果
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --progress=false
```

### 4. 多行格式输出
```bash
# 多行输出（适合人类阅读）
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --format strm

# 静默多行输出可用镜像
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --quiet \
  --format strm \
  --print-result=false \
  --print-advice=false
```

### 5. 组合使用场景

#### 测试 + 更新 + 显示失败列表
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --update \
  --print-fails \
  --format strm
```

#### 静默更新（适合自动化脚本）
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --quiet \
  --update \
  --progress=false
```

#### 完整输出控制
```bash
sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --print-result=true \
  --print-advice=true \
  --print-fails=true \
  --progress=true \
  --format strm
```

### 6. 配合其他工具使用

#### 提取可用镜像列表到环境变量
```bash
AVAILABLE=$(sudo bin/dogo \
  --file /etc/docker/daemon.json \
  --quiet \
  --format stro \
  --print-result=false \
  --print-advice=false)
echo "可用镜像: $AVAILABLE"
```

#### 批量处理多个配置文件
```bash
for config in /etc/docker/daemon.json ./config.json; do
  sudo bin/dogo --file "$config" --quiet --update
done
```

#### 监控镜像可用性（定时任务）
```bash
#!/bin/bash
# 每天凌晨2点测试镜像可用性
0 2 * * * /usr/local/bin/dogo --file /etc/docker/daemon.json --quiet --update --print-fails
```

## 📁 支持的JSON格式

### 标准 Docker daemon.json 格式
```json
{
  "registry-mirrors": [
    "https://mirror.aliyuncs.com",
    "https://docker.m.daocloud.io",
    "https://docker.nju.edu.cn"
  ],
  "insecure-registries": [
    "docker.mirrors.ustc.edu.cn",
    "hub-mirror.c.163.com"
  ],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
```

### 自定义配置格式
```json
{
  "mirrors": [
    "https://mirror.aliyuncs.com",
    "https://docker.m.daocloud.io"
  ],
  "timeout": 10,
  "concurrency": 20,
  "print_result": true,
  "print_advice": true,
  "print_fails": false,
  "show_progress": true,
  "quiet": false,
  "format": "stro",
  "update": false
}
```

## 📊 输出格式说明

### 1. 单行格式 (stro)
```
https://mirror.aliyuncs.com,https://docker.m.daocloud.io,https://docker.nju.edu.cn
```
- 适合：脚本解析、环境变量、批量处理

### 2. 多行格式 (strm)
```
https://mirror.aliyuncs.com
https://docker.m.daocloud.io
https://docker.nju.edu.cn
```
- 适合：人类阅读、配置文件生成

### 3. 完整输出示例
```
🚀 开始测试 18 个镜像加速器...
📁 配置文件: /etc/docker/daemon.json
⏱️  超时: 10s, 并发数: 20
📊 进度条已启用

进度: [██████████████████████████████████████████████████] 100.0% (18/18)

📊 测试结果:
--------------------------------------------------------------------------------
✅ https://mirror.aliyuncs.com                      延迟: 234ms
✅ https://docker.m.daocloud.io                    延迟: 312ms
❌ https://docker.registry.cyou                    所有端点测试失败
✅ https://docker.nju.edu.cn                        延迟: 456ms
--------------------------------------------------------------------------------
📈 统计: 12/18 可用, 平均延迟: 342ms

⚡ 最快镜像 (Top 3):
  1. https://mirror.aliyuncs.com (延迟: 234ms)
  2. https://docker.m.daocloud.io (延迟: 312ms)
  3. https://docker.nju.edu.cn (延迟: 456ms)

💡 优化建议:
--------------------------------------------------------------------------------
推荐使用的镜像加速器（按速度排序）:
  1. https://mirror.aliyuncs.com (延迟: 234ms)
  2. https://docker.m.daocloud.io (延迟: 312ms)
  3. https://docker.nju.edu.cn (延迟: 456ms)

📋 可用镜像列表:
--------------------------------------------------------------------------------
https://mirror.aliyuncs.com,https://docker.m.daocloud.io,https://docker.nju.edu.cn

❌ 不可用镜像列表:
--------------------------------------------------------------------------------
https://docker.registry.cyou,https://docker-cf.registry.cyou

📊 不可用镜像: 2 个
```

## 🐳 容器化部署（Docker）

### 构建镜像
```bash
docker build --progress=plain -f Dockerfile.dogo --target runtime -t dogo:latest .
```

### 运行测试
```bash
# 挂载宿主机 docker 配置文件到容器内
docker run -it --rm -v /etc/docker:/etc/docker:ro dogo:latest dogo --file /etc/docker/daemon.json --update
```

## 📦 编译安装

### 本地编译
```bash
# 编译二进制文件
go build -o bin/dogo ./cmd/dogo/main.go

# 静态编译（适合容器）
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -ldflags="-w -s -extldflags=-static" \
  -trimpath \
  -o bin/dogo ./cmd/dogo/main.go
```

### 使用 Makefile
```makefile
.PHONY: build clean install

build:
	go build -o bin/dogo ./cmd/dogo/main.go

install:
	go install github.com/ymc-github/dogo/cmd/dogo@latest

clean:
	rm -rf bin/
```

## ⚠️ 注意事项

1. **系统配置权限**：`/etc/docker/daemon.json` 属于系统保护文件，读写必须加 `sudo`
2. **文件格式要求**：配置文件必须包含 `registry-mirrors` 或 `mirrors` 字段
3. **网络环境**：测试需要访问外网镜像源，确保网络连通性
4. **并发控制**：并发数过高可能触发限流，建议保持在 10-30 之间
5. **超时设置**：网络慢的环境可适当增加超时时间
6. **备份建议**：更新配置文件前建议先备份原文件
7. **Docker重载**：修改配置后需要重启 Docker 服务生效

## 🛡️ 特性亮点

- ✅ **零配置开箱即用**：仅需 `--file` 一个必填参数即可运行
- ✅ **智能测试**：并发测试、自动排序、去重处理
- ✅ **运维友好**：静默模式、多种输出格式、完善日志
- ✅ **安全可控**：支持预览、可控制是否更新文件
- ✅ **容器友好**：静态编译，无额外依赖
- ✅ **自动化适配**：适合 CI/CD、定时任务、批量脚本
- ✅ **灵活输出**：支持人类可读和机器解析两种格式
- ✅ **性能优化**：并发测试、快速响应

## 🔧 故障排查

### 常见问题

#### 1. 配置文件不存在
```
❌ 错误: 配置文件不存在: /etc/docker/daemon.json
```
**解决方案**：确认文件路径正确，使用 `--file` 指定正确路径

#### 2. 没有找到镜像列表
```
❌ 配置文件中没有找到 registry-mirrors 或 mirrors 字段
```
**解决方案**：检查配置文件是否包含 `registry-mirrors` 或 `mirrors` 字段

#### 3. 权限不足
```
❌ 更新配置文件失败: permission denied
```
**解决方案**：使用 `sudo` 执行命令

#### 4. 所有镜像测试失败
```
⚠️  没有可用的镜像加速器
```
**解决方案**：检查网络连接，尝试增加超时时间 `--timeout 20`

### 调试模式
```bash
# 开启详细日志（需要添加 -v 参数支持）
dogo --file daemon.json -v
```

## 📝 更新日志

### v1.0.0
- ✅ 初始版本发布
- ✅ 支持批量镜像测试
- ✅ 支持自动更新配置
- ✅ 支持多种输出格式
- ✅ 支持静默模式
- ✅ 支持进度条控制
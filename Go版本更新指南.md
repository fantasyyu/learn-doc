# 🔄 Go 版本更新指南（Windows）

> 本机工作目录：`D:\aiworkspace\go_learn\go学习`

---

## 🔍 当前环境状态

- `go version` 命令：**未找到**（环境变量可能未配置或旧版本未正确安装）
- `C:\Go` 目录存在，但可能不完整
- 结论：**建议直接重新安装最新版 Go**

---

## ✅ 推荐方法：官网下载覆盖安装（最简单）

### Step 1：下载最新版 Go

1. 打开 [https://go.dev/dl/](https://go.dev/dl/)
2. 找到 **Microsoft Windows** 部分
3. 下载：`go1.23.x.windows-amd64.msi`（当前最新稳定版）
   - 或根据你的系统选择 `go1.23.x.windows-arm64.msi`（ARM 设备）

### Step 2：卸载旧版本（可选但推荐）

**方式 A：控制面板卸载**
1. 打开「控制面板」→「程序和功能」
2. 找到 **Go Programming Language**，右键「卸载」

**方式 B：手动删除**
```powershell
# 如果之前装在 C:\Go，直接删除该目录
Remove-Item -Recurse -Force "C:\Go"
```

### Step 3：安装新版本

1. 双击下载的 `.msi` 文件
2. 一路「Next」，建议保持默认安装路径：`C:\Go`
3. 安装程序会自动配置环境变量

### Step 4：验证安装

打开一个新的 **PowerShell** 或 **CMD** 窗口（**必须新开窗口**），执行：

```powershell
go version
```

期望输出类似：
```
go version go1.23.0 windows/amd64
```

### Step 5：配置 GOPATH（推荐）

Go 1.8+ 默认使用 `~/go` 作为 GOPATH，但建议显式设置：

1. 右键「此电脑」→「属性」→「高级系统设置」→「环境变量」
2. 在「用户变量」中点击「新建」：
   - 变量名：`GOPATH`
   - 变量值：`D:\aiworkspace\go_learn\go学习\gopath`
3. 编辑「用户变量」中的 `Path`，添加：
   - `%GOPATH%\bin`
4. 点击「确定」保存

### Step 6：初始化学习项目

```powershell
cd "D:\aiworkspace\go_learn\go学习"
go mod init go学习
go version
```

---

## 🛠️ 备选方法：命令行包管理器安装

### 方法 A：使用 Scoop（推荐，方便后续升级）

```powershell
# 安装 Scoop（如未安装）
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression

# 安装 Go
scoop install go

# 以后升级只需
scoop update go
```

### 方法 B：使用 Winget（Windows 自带）

```powershell
# 安装 Go
winget install GoLang.Go

# 以后升级
winget upgrade GoLang.Go
```

### 方法 C：使用 Chocolatey

```powershell
# 安装 Go
choco install golang

# 以后升级
choco upgrade golang
```

---

## 🧹 常见问题排查

### 问题 1：安装后 `go version` 仍显示旧版本

**原因**：旧版本的 Go 路径仍在环境变量中排在新版本前面。

**解决**：
```powershell
# 查看 go 实际执行路径
where go

# 如果显示多个路径，调整环境变量 Path 的顺序
# 确保 C:\Go\bin 排在最前面
```

### 问题 2：`go` 命令找不到

**原因**：环境变量未生效或未配置。

**解决**：
```powershell
# 临时添加（仅当前窗口有效）
$env:Path += ";C:\Go\bin"

# 永久添加（需要管理员权限）
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\Go\bin", "User")
```

### 问题 3：VSCode 找不到 Go

**解决**：
1. 重启 VSCode
2. 按 `Ctrl+Shift+P` → `Go: Install/Update Tools`
3. 全选 → 安装

### 问题 4：代理问题导致 `go get` 失败

**解决**：配置国内代理
```powershell
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn
```

---

## 📁 本机目录规划建议

```
D:\aiworkspace\go_learn\go学习\          # 学习项目根目录
├── 学习计划.md                           # 本计划
├── hello.go                              # 第一个程序
├── go.mod                                # 模块文件
├── stage1/                               # 阶段1练习
│   ├── calc.go
│   └── guess_number.go
├── stage2/                               # 阶段2练习
│   └── library_system.go
├── stage3/                               # 阶段3练习
│   └── http_server.go
├── stage4/                               # 阶段4练习
│   └── concurrent_downloader.go
├── stage5/                               # 阶段5练习
│   └── gin_demo/
├── stage6/                               # 阶段6实战
│   └── todo_api/
└── stage7/                               # 阶段7项目
    └── blog_backend/
```

---

## ✅ 安装完成后检查清单

```
□ go version 显示最新版本号
□ go env GOPATH 显示正确路径
□ VSCode Go 插件已安装且工具已更新
□ go mod init 能正常创建 go.mod
□ go run hello.go 能正常运行
□ go get 命令能正常下载依赖（已配置代理）
```

---

> 💡 如果在安装或更新过程中遇到任何问题，把错误信息贴出来，我帮你解决！

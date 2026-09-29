package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	fmt.Println("=== Go Modules 核心概念 ===")
	fmt.Println()
	fmt.Println("1. go.mod — 模块定义文件")
	fmt.Println("   module 语句: 模块路径（通常是仓库地址）")
	fmt.Println("   require 语句: 依赖列表及版本")
	fmt.Println("   replace 语句: 替换依赖（本地调试常用）")
	fmt.Println()
	fmt.Println("2. go.sum — 依赖校验文件")
	fmt.Println("   记录每个依赖的哈希值，确保下载的代码没被篡改")
	fmt.Println("   不要手动编辑，提交到版本控制")
	fmt.Println()
	fmt.Println("3. 语义化版本")
	fmt.Println("   v1.2.3  = 主版本.次版本.补丁版本")
	fmt.Println("   v0.x.x  = 不稳定版本，API 可能 breaking change")
	fmt.Println("   /v2+    = 主版本 >= 2 时，模块路径必须带 /v2 后缀")
	fmt.Println()
	fmt.Println("4. 常用命令")
	fmt.Println("   go mod init <模块名>   初始化模块")
	fmt.Println("   go get <包>            下载依赖（加到 require）")
	fmt.Println("   go get -u              升级所有依赖")
	fmt.Println("   go mod tidy            清理未使用的依赖（最常用！）")
	fmt.Println("   go mod download        下载所有依赖到本地缓存")
	fmt.Println("   go list -m all         查看所有依赖版本")
	fmt.Println("   go mod why <包>        查看为什么需要这个依赖")
	fmt.Println()
	fmt.Println("=== 本模块信息 ===")
	fmt.Println("模块名: example.com/day29")
	fmt.Println("Go版本: 1.22")
	fmt.Println("依赖见 go.mod 文件（下面下载了一个真实依赖）")
	// ============================================================
	//  5. 实际使用第三方依赖（go get 下载的 uuid 包）
	// ============================================================
	fmt.Println("=== 使用依赖: github.com/google/uuid ===")
	for i := 0; i < 3; i++ {
		id := uuid.New()
		fmt.Printf("UUID: %s\n", id.String())
	}
}

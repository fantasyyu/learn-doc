package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// ============================================================
	//  1. flag 包 — 解析命令行参数
	// ============================================================

	// 定义参数
	name := flag.String("name", "World", "你的名字")
	age := flag.Int("age", 18, "你的年龄")
	verbose := flag.Bool("v", false, "显示详细信息")

	// 解析（必须调用！）
	flag.Parse()

	fmt.Println("=== flag 命令行参数 ===")
	fmt.Printf("你好, %s! 年龄: %d\n", *name, *age)
	if *verbose {
		fmt.Println("详细模式已开启")
	}

	// 获取非 flag 参数
	args := flag.Args()
	if len(args) > 0 {
		fmt.Printf("额外参数: %v\n", args)
	}

	// ============================================================
	//  2. log 包 — 日志记录
	// ============================================================
	fmt.Println("\n=== log 日志 ===")

	// 默认 logger 输出到 stderr，带时间戳
	log.Println("这是一条普通日志")
	log.Printf("格式化日志: %s", "hello")

	// 设置前缀和标志
	log.SetPrefix("[APP] ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("带前缀和文件位置的日志")

	// 自定义 logger
	fmt.Println("\n=== 自定义 logger ===")
	logger := log.New(os.Stdout, "[CUSTOM] ", log.LstdFlags|log.Lmicroseconds)
	logger.Println("自定义 logger 的消息")

	// 写入文件
	f, _ := os.Create("phase3/day19/app.log")
	defer f.Close()
	fileLogger := log.New(f, "[FILE] ", log.LstdFlags)
	fileLogger.Println("这条日志写入了文件")
	fmt.Println("✓ 日志已写入 phase3/day19/app.log")

	// Fatal / Panic
	// log.Fatal("致命错误，程序退出")   // os.Exit(1)
	// log.Panic("恐慌错误")             // panic()
}

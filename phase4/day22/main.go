package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	// ============================================================
	//  1. goroutine 基础 — go 关键字启动协程
	// ============================================================
	fmt.Println("=== 1. goroutine 基础 ===")

	go fmt.Println("我是一个 goroutine") // 不等待结果，继续往下执行
	fmt.Println("主函数继续执行")
	time.Sleep(100 * time.Millisecond) // 等待 goroutine 完成

	// ============================================================
	//  2. goroutine 是轻量级的
	// ============================================================
	fmt.Println("\n=== 2. 启动大量 goroutine ===")

	var wg sync.WaitGroup // 等待组：等所有 goroutine 完成
	start := time.Now()

	for i := 0; i < 10; i++ {
		wg.Add(1) // 计数 +1
		go func(id int) {
			defer wg.Done() // 完成时计数 -1
			fmt.Printf("worker-%d 启动\n", id)
			time.Sleep(50 * time.Millisecond)
		}(i)
	}

	wg.Wait() // 等所有 goroutine 完成
	fmt.Printf("10 个 goroutine 耗时: %v\n", time.Since(start))

	fmt.Printf("\n当前 goroutine 数量: %d\n", runtime.NumGoroutine())
	fmt.Printf("CPU 核数: %d\n", runtime.NumCPU())

	// ============================================================
	//  3. 闭包陷阱 — goroutine 捕获循环变量
	// ============================================================
	fmt.Println("\n=== 3. 闭包陷阱（循环变量共享）===")

	// ❌ 错误：所有 goroutine 共享同一个变量 i
	for i := 0; i < 3; i++ {
		go func() {
			// 这里捕获的是外部循环变量的引用
			fmt.Printf("错误方式: goroutine 看到 i=%d\n", i)
		}()
	}
	time.Sleep(200 * time.Millisecond)

	// ✅ 正确：通过参数传递，每个 goroutine 有自己的副本
	for i := 0; i < 3; i++ {
		go func(id int) {
			fmt.Printf("正确方式: goroutine 看到 id=%d\n", id)
		}(i)
	}
	time.Sleep(200 * time.Millisecond)

	// ============================================================
	//  4. goroutine 之间通过 channel 通信（下一课细讲）
	// ============================================================
	fmt.Println("\n=== 4. channel 初体验 ===")

	ch := make(chan string) // 创建无缓冲 channel

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch <- "goroutine 完成了！" // 发送数据
	}()

	msg := <-ch // 接收数据（阻塞等待）
	fmt.Println("收到:", msg)

	fmt.Println("\n=== 总结 ===")
	fmt.Println("1. go 关键字启动 goroutine，是非阻塞的")
	fmt.Println("2. sync.WaitGroup 等待所有 goroutine 完成")
	fmt.Println("3. 循环变量陷阱：用参数传递，不要捕获循环变量")
	fmt.Println("4. channel 是 goroutine 之间通信的桥梁")
}

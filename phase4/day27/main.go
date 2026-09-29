package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// ============================================================
	//  1. 原子操作 — 比 Mutex 更轻量
	// ============================================================
	fmt.Println("=== 1. 原子操作 vs Mutex ===")

	// 场景：简单计数器
	// 方式A：Mutex
	var mu sync.Mutex
	var counterA int64
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			counterA++
			mu.Unlock()
		}()
	}
	wg.Wait()
	fmt.Printf("Mutex:  %d, 耗时: %v\n", counterA, time.Since(start))

	// 方式B：atomic（无锁，更快）
	var counterB int64
	start = time.Now()
	for i := 0; i < 100000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counterB, 1)
		}()
	}
	wg.Wait()
	fmt.Printf("atomic: %d, 耗时: %v\n", counterB, time.Since(start))

	// ============================================================
	//  2. atomic 常用操作
	// ============================================================
	fmt.Println("\n=== 2. atomic 常用操作 ===")

	var val int64 = 100

	// 读取
	fmt.Printf("Load: %d\n", atomic.LoadInt64(&val))

	// 存储
	atomic.StoreInt64(&val, 200)
	fmt.Printf("Store 200 后: %d\n", atomic.LoadInt64(&val))

	// CAS（Compare-And-Swap）— 乐观锁的核心
	// 如果当前值等于 old，则替换为 new，返回 true；否则返回 false
	swapped := atomic.CompareAndSwapInt64(&val, 200, 300)
	fmt.Printf("CAS(200→300): %v, 当前值: %d\n", swapped, atomic.LoadInt64(&val))

	swapped = atomic.CompareAndSwapInt64(&val, 200, 999)
	fmt.Printf("CAS(200→999): %v, 当前值: %d ← 失败，值没变\n", swapped, atomic.LoadInt64(&val))

	// Add（负数就是减）
	atomic.AddInt64(&val, -50)
	fmt.Printf("Add(-50) 后: %d\n", atomic.LoadInt64(&val))

	// ============================================================
	//  3. atomic.Value — 原子存取任意值
	// ============================================================
	fmt.Println("\n=== 3. atomic.Value ===")

	var config atomic.Value
	config.Store(map[string]string{"env": "dev", "version": "1.0"})

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			m := config.Load().(map[string]string)
			fmt.Printf("读取配置: %v\n", m)
			time.Sleep(10 * time.Millisecond)
		}
	}()
	go func() {
		defer wg.Done()
		time.Sleep(20 * time.Millisecond)
		config.Store(map[string]string{"env": "prod", "version": "2.0"})
		fmt.Println("配置已更新为 prod")
	}()
	wg.Wait()

	// ============================================================
	//  4. 什么时候用什么？
	// ============================================================
	fmt.Println("\n=== 4. 选型建议 ===")
	fmt.Println("atomic  — 简单计数器、标志位、配置读取（高性能）")
	fmt.Println("Mutex   — 保护复杂临界区（多变量操作、数据结构）")
	fmt.Println("channel — goroutine 之间传递数据（Go 哲学首选）")
	fmt.Println("")
	fmt.Println("性能: atomic > Mutex > channel（但可读性相反）")

	// ============================================================
	//  5. ⭐ 竞态检测 — go run -race
	// ============================================================
	fmt.Println("\n=== 5. 竞态检测 ===")

	// 这段代码有竞态问题！
	// 运行: go run -race phase4/day27/main.go
	// 会输出: WARNING: DATA RACE
	racyCounter := 0
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			racyCounter++ // ❌ 与其他 goroutine 竞争
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			racyCounter++ // ❌ 与其他 goroutine 竞争
		}
	}()
	wg.Wait()
	fmt.Printf("竞态计数（结果随机）: %d\n", racyCounter)
	fmt.Println("\n👉 试试: go run -race phase4/day27/main.go")
}

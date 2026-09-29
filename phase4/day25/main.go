package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// ============================================================
	//  1. 竞态条件演示 — 不加锁的错误
	// ============================================================
	fmt.Println("=== 1. 竞态条件（错误示范）===")

	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // ❌ 非线程安全！多个 goroutine 同时读写
		}()
	}
	wg.Wait()
	fmt.Printf("期望: 1000, 实际: %d ← 每次运行结果可能不同！\n", counter)

	// ============================================================
	//  2. Mutex — 互斥锁，保护共享变量
	// ============================================================
	fmt.Println("\n=== 2. Mutex 保护共享变量 ===")

	var mu sync.Mutex
	safeCounter := 0

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()   // 加锁：同一时间只有一个 goroutine 能进入
			safeCounter++ // 临界区
			mu.Unlock() // 解锁
		}()
	}
	wg.Wait()
	fmt.Printf("期望: 1000, 实际: %d ✓\n", safeCounter)

	// ============================================================
	//  3. RWMutex — 读写锁（读多写少场景）
	// ============================================================
	fmt.Println("\n=== 3. RWMutex 读写锁 ===")

	var rwMu sync.RWMutex
	data := make(map[string]int)

	// 写入者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 5; i++ {
			rwMu.Lock() // 写锁
			data[fmt.Sprintf("key%d", i)] = i
			fmt.Printf("写入: key%d=%d\n", i, i)
			rwMu.Unlock()
			time.Sleep(50 * time.Millisecond)
		}
	}()

	// 多个读取者
	for r := 0; r < 3; r++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < 5; i++ {
				rwMu.RLock() // 读锁（多个读可以同时进行）
				val, ok := data[fmt.Sprintf("key%d", i)]
				fmt.Printf("  读者%d 读取: key%d=%d (ok=%v)\n", id, i, val, ok)
				rwMu.RUnlock()
				time.Sleep(20 * time.Millisecond)
			}
		}(r)
	}
	wg.Wait()

	// ============================================================
	//  4. Once — 只执行一次
	// ============================================================
	fmt.Println("\n=== 4. sync.Once ===")

	var once sync.Once
	initValue := 0

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			once.Do(func() { // 只执行一次，即使多个 goroutine 同时调用
				fmt.Println("初始化（只执行一次）")
				initValue = 100
			})
		}()
	}
	wg.Wait()
	fmt.Printf("initValue = %d\n", initValue)

	// ============================================================
	//  5. sync.Map — 并发安全的 map
	// ============================================================
	fmt.Println("\n=== 5. sync.Map ===")

	var sm sync.Map

	// 并发写入
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			sm.Store(fmt.Sprintf("user%d", id), id*10)
		}(i)
	}
	wg.Wait()

	// 读取
	sm.Range(func(key, value interface{}) bool {
		fmt.Printf("%s: %d\n", key, value)
		return true // 继续遍历
	})

	val, ok := sm.Load("user5")
	fmt.Printf("\nuser5: %v, ok=%v\n", val, ok)

	// ============================================================
	//  6. Pool — 对象复用池
	// ============================================================
	fmt.Println("\n=== 6. sync.Pool（对象复用）===")

	var pool = sync.Pool{
		New: func() interface{} { // 池为空时创建新对象
			fmt.Println("创建新对象")
			return make([]byte, 1024)
		},
	}

	// 获取对象
	buf := pool.Get().([]byte)
	fmt.Printf("获取到 buffer: %d bytes\n", len(buf))

	// 使用完归还
	pool.Put(buf)
	fmt.Println("对象已归还到池中")

	fmt.Println("\n=== 总结 ===")
	fmt.Println("Mutex  — 互斥锁，保护临界区（写操作）")
	fmt.Println("RWMutex — 读写锁（读多写少）")
	fmt.Println("Once   — 确保某操作只执行一次（如初始化）")
	fmt.Println("Map    — 并发安全的 map")
	fmt.Println("Pool   — 对象复用，减少 GC 压力")
}

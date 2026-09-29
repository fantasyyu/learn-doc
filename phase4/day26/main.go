package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 模拟一个耗时的子任务
func worker(ctx context.Context, id int, duration time.Duration) error {
	select {
	case <-time.After(duration):
		fmt.Printf("worker-%d 完成（耗时 %v）\n", id, duration)
		return nil
	case <-ctx.Done():
		fmt.Printf("worker-%d 被取消\n", id)
		return ctx.Err()
	}
}

func main() {
	// ============================================================
	//  1. 一个 context 取消所有 goroutine
	// ============================================================
	fmt.Println("=== 1. context 取消所有 goroutine ===")

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	tasks := []time.Duration{
		200 * time.Millisecond,
		400 * time.Millisecond,
		600 * time.Millisecond,
	}

	for i, d := range tasks {
		wg.Add(1)
		go func(id int, dur time.Duration) {
			defer wg.Done()
			worker(ctx, id, dur)
		}(i, d)
	}

	// 300ms 后取消所有任务
	go func() {
		time.Sleep(300 * time.Millisecond)
		fmt.Println("→ 发送取消信号（300ms 到了）")
		cancel()
	}()

	wg.Wait()
	fmt.Println("所有 worker 已结束")

	// ============================================================
	//  2. 并发搜索 — 找到第一个结果就取消其他
	// ============================================================
	fmt.Println("\n=== 2. 并发搜索（找到结果即取消其他）===")

	searchEngine := func(ctx context.Context, engine string, delay time.Duration, result chan<- string) {
		select {
		case <-time.After(delay):
			select {
			case result <- fmt.Sprintf("%s 的结果", engine):
			case <-ctx.Done(): // 结果没人接收（可能已有其他引擎返回了）
				fmt.Printf("%s 的结果被丢弃\n", engine)
			}
		case <-ctx.Done():
			fmt.Printf("%s 被取消了\n", engine)
		}
	}

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	results := make(chan string, 3)

	// 3 个搜索引擎并发搜索
	go searchEngine(ctx2, "百度", 100*time.Millisecond, results)
	go searchEngine(ctx2, "Google", 200*time.Millisecond, results)
	go searchEngine(ctx2, "Bing", 300*time.Millisecond, results)

	// 取第一个结果
	first := <-results
	fmt.Printf("✓ 使用: %s\n", first)
	cancel2() // 取消其他搜索

	time.Sleep(400 * time.Millisecond) // 等其他 goroutine 处理取消

	// ============================================================
	//  3. context 传值 — 在 goroutine 之间传递请求信息
	// ============================================================
	fmt.Println("\n=== 3. context 传值 ===")

	type ctxKey string // 避免 key 冲突的最佳实践：自定义类型

	var process func(ctx context.Context, depth int)
	process = func(ctx context.Context, depth int) {
		requestID := ctx.Value(ctxKey("request_id"))
		fmt.Printf("第%d层处理: request_id=%v\n", depth, requestID)

		if depth < 3 {
			// 值沿着调用链传递
			process(ctx, depth+1)
		}
		fmt.Printf("第%d层处理结束: request_id=%v\n", depth, requestID)
	}

	ctx3 := context.WithValue(context.Background(), ctxKey("request_id"), "req-2024")
	process(ctx3, 1)

	// ============================================================
	//  4. 超时控制 — 整个操作不能超过 X 秒
	// ============================================================
	fmt.Println("\n=== 4. 整体超时控制 ===")

	ctx4, cancel4 := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel4()

	start := time.Now()
	var wg4 sync.WaitGroup

	// 5 个并发任务，每个随机耗时 100-800ms
	for i := 0; i < 5; i++ {
		wg4.Add(1)
		go func(id int) {
			defer wg4.Done()
			d := time.Duration(rand.Intn(700)+100) * time.Millisecond
			worker(ctx4, id, d)
		}(i)
	}

	wg4.Wait()
	fmt.Printf("总耗时: %v（超时限制 500ms）\n", time.Since(start))

	fmt.Println("\n=== 总结 ===")
	fmt.Println("1. context.WithCancel → cancel() 同时取消所有子 goroutine")
	fmt.Println("2. 并发搜索：找到第一个结果后 cancel() 取消其他")
	fmt.Println("3. WithValue 传值：用自定义类型作为 key，避免冲突")
	fmt.Println("4. WithTimeout 控制整个并发操作的超时时间")
}

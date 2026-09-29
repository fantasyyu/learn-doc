package main

import (
	"fmt"
	"time"
)

func main() {
	// ============================================================
	//  1. 无缓冲 channel — 同步通信
	//  发送和接收必须同时准备好，否则阻塞
	// ============================================================
	fmt.Println("=== 1. 无缓冲 channel（同步）===")

	ch1 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		fmt.Println("发送前...")
		ch1 <- "消息" // 接收方没准备好前，这里会一直阻塞
		fmt.Println("发送后...") // 这行不会立即执行
	}()

	time.Sleep(200 * time.Millisecond)
	msg := <-ch1 // 接收（发送方还没发时会阻塞）
	fmt.Println("收到:", msg)

	// ============================================================
	//  2. 有缓冲 channel — 异步通信
	// ============================================================
	fmt.Println("\n=== 2. 有缓冲 channel（异步）===")

	ch2 := make(chan int, 3) // 缓冲大小 3

	ch2 <- 1 // 不阻塞，缓冲没满
	ch2 <- 2
	ch2 <- 3
	fmt.Println("已发送 3 个值，不阻塞")

	go func() {
		time.Sleep(100 * time.Millisecond)
		for i := 4; i <= 6; i++ {
			ch2 <- i // 缓冲满了会阻塞
		}
		close(ch2) // 关闭 channel
	}()

	// range 遍历 channel，直到关闭
	for v := range ch2 {
		fmt.Printf("接收: %d\n", v)
	}

	// ============================================================
	//  3. 单向 channel — 限制方向
	// ============================================================
	fmt.Println("\n=== 3. 单向 channel ===")

	// 生产者：只能发送
	producer := func(out chan<- int) { // 只写
		for i := 1; i <= 3; i++ {
			out <- i
		}
		close(out)
	}

	// 消费者：只能接收
	consumer := func(in <-chan int) { // 只读
		for v := range in {
			fmt.Printf("消费: %d\n", v)
		}
	}

	ch3 := make(chan int)
	go producer(ch3)
	consumer(ch3)

	// ============================================================
	//  4. 关闭 channel 的规则
	// ============================================================
	fmt.Println("\n=== 4. close 规则 ===")

	ch4 := make(chan int, 2)
	ch4 <- 1
	ch4 <- 2
	close(ch4)

	// 从已关闭的 channel 接收：返回零值，ok=false
	v1, ok1 := <-ch4
	v2, ok2 := <-ch4
	v3, ok3 := <-ch4
	fmt.Printf("接收1: %d, ok=%v\n", v1, ok1) // 1, true
	fmt.Printf("接收2: %d, ok=%v\n", v2, ok2) // 2, true
	fmt.Printf("接收3: %d, ok=%v ← 已关闭，返回零值\n", v3, ok3) // 0, false

	fmt.Println("\n规则:")
	fmt.Println("  只有发送方可以 close channel")
	fmt.Println("  向已关闭的 channel 发送数据 → panic！")
	fmt.Println("  从已关闭的 channel 接收 → 返回零值，ok=false")

	// ============================================================
	//  5. 生产者消费者模式
	// ============================================================
	fmt.Println("\n=== 5. 生产者消费者模式 ===")

	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// 3 个 worker
	for w := 1; w <= 3; w++ {
		go func(id int) {
			for j := range jobs { // 从 jobs 取任务
				time.Sleep(50 * time.Millisecond) // 模拟耗时
				results <- j * 2 // 写入结果
				fmt.Printf("worker-%d 处理 job-%d\n", id, j)
			}
		}(w)
	}

	// 发送 5 个任务
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// 收集结果
	for i := 0; i < 5; i++ {
		fmt.Printf("结果: %d\n", <-results)
	}
}

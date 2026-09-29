package main

import (
	"fmt"
	"time"
)

func main() {
	// ============================================================
	//  1. select 基础 — 同时监听多个 channel
	// ============================================================
	fmt.Println("=== 1. select 基础 ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "来自 ch1"
	}()
	go func() {
		time.Sleep(50 * time.Millisecond)
		ch2 <- "来自 ch2"
	}()

	// select：哪个 channel 先有数据就执行哪个 case
	select {
	case msg1 := <-ch1:
		fmt.Println(msg1)
	case msg2 := <-ch2:
		fmt.Println(msg2) // ch2 先就绪，执行这个
	}

	// ============================================================
	//  2. select + 超时 — 防止永久阻塞
	// ============================================================
	fmt.Println("\n=== 2. select + 超时 ===")

	ch3 := make(chan string) // 永远不会收到数据
	// go func() {
	// 	time.Sleep(50 * time.Millisecond)
	// 	ch3 <- "来自 ch3"
	// }()
	select {
	case msg := <-ch3:
		fmt.Println(msg)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("⏰ 超时了！")
	}

	// ============================================================
	//  3. select + default — 非阻塞操作
	// ============================================================
	fmt.Println("\n=== 3. select + default（非阻塞）===")

	ch4 := make(chan int)

	select {
	case v := <-ch4:
		fmt.Println(v)
	default:
		fmt.Println("channel 是空的，不阻塞！")
	}

	// 非阻塞发送
	select {
	case ch4 <- 1:
		fmt.Println("发送成功")
	default:
		fmt.Println("channel 满了或没有接收者，不阻塞！")
	}

	// ============================================================
	//  4. 心跳模式 — 定时执行任务
	// ============================================================
	fmt.Println("\n=== 4. 心跳模式 ===")

	heartbeat := time.NewTicker(200 * time.Millisecond)
	defer heartbeat.Stop()

	timeout := time.After(1 * time.Second)

	count := 0

	for {
		select {
		case t := <-heartbeat.C:
			count++
			fmt.Printf("心跳 #%d: %v\n", count, t.Format("15:04:05.000"))
		case <-timeout:
			fmt.Println("⏰ 心跳结束")
			goto done
		}
	}
done:
	fmt.Println("心跳模式结束")
	// ============================================================
	//  5. 多 channel 随机选择 — 当多个 case 同时就绪
	// ============================================================
	fmt.Println("\n=== 5. 多个 case 同时就绪（随机选择）===")

	a := make(chan string, 1)
	b := make(chan string, 1)
	a <- "A"
	b <- "B"

	for i := 0; i < 2; i++ {
		select {
		case msg := <-a:
			fmt.Println("选了:", msg)
		case msg := <-b:
			fmt.Println("选了:", msg)
		}
	}

	// ============================================================
	//  6. 实际应用：并发请求，谁快用谁
	// ============================================================
	fmt.Println("\n=== 6. 并发请求，取最快的结果 ===")

	fastServer := make(chan string)
	slowServer := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		fastServer <- "服务器A (100ms)"
	}()
	go func() {
		time.Sleep(500 * time.Millisecond)
		slowServer <- "服务器B (500ms)"
	}()

	select {
	case result := <-fastServer:
		fmt.Println("✓ 使用:", result)
	case result := <-slowServer:
		fmt.Println("✓ 使用:", result)
	case <-time.After(300 * time.Millisecond):
		fmt.Println("⏰ 请求超时")
	}

	fmt.Println("\n=== 总结 ===")
	fmt.Println("1. select 同时监听多个 channel，谁就绪执行谁")
	fmt.Println("2. time.After + select = 超时控制")
	fmt.Println("3. default = 非阻塞操作")
	fmt.Println("4. time.Ticker + select = 定时任务/心跳")
	fmt.Println("5. 多个 case 同时就绪 → 随机选择")
}

package main

import (
	"context"
	"fmt"
	"time"
)

// 模拟一个耗时操作
func slowTask(ctx context.Context, name string, duration time.Duration) error {
	select {
	case <-time.After(duration): // 任务完成
		fmt.Printf("✓ %s 完成\n", name)
		return nil
	case <-ctx.Done(): // 收到取消信号
		return fmt.Errorf("%s 被取消: %v", name, ctx.Err())
	}
}

func main() {
	// ============================================================
	//  1. context.Background() — 根 context
	// ============================================================
	fmt.Println("=== context 基础 ===")
	ctx := context.Background()
	fmt.Printf("Background: %v\n", ctx)

	// ============================================================
	//  2. WithCancel — 手动取消
	// ============================================================
	fmt.Println("\n=== WithCancel 手动取消 ===")
	ctx2, cancel := context.WithCancel(ctx)

	go func() {
		time.Sleep(500 * time.Millisecond)
		fmt.Println("→ 发送取消信号")
		cancel() // 取消！
	}()

	err := slowTask(ctx2, "任务A", 2*time.Second)
	if err != nil {
		fmt.Println("结果:", err)
	}

	// ============================================================
	//  3. WithTimeout — 超时自动取消
	// ============================================================
	fmt.Println("\n=== WithTimeout 超时控制 ===")
	ctx3, cancel3 := context.WithTimeout(ctx, 1*time.Second)
	defer cancel3() // 好习惯：即使超时也会触发，避免资源泄漏

	err = slowTask(ctx3, "任务B(需2秒)", 2*time.Second)
	if err != nil {
		fmt.Println("结果:", err)
	}

	// ============================================================
	//  4. WithDeadline — 截止时间
	// ============================================================
	fmt.Println("\n=== WithDeadline 截止时间 ===")
	deadline := time.Now().Add(800 * time.Millisecond)
	ctx4, cancel4 := context.WithDeadline(ctx, deadline)
	defer cancel4()

	err = slowTask(ctx4, "任务C", 2*time.Second)
	if err != nil {
		fmt.Println("结果:", err)
	}

	// ============================================================
	//  5. WithValue — 传递键值对（如 requestID、userID）
	// ============================================================
	fmt.Println("\n=== WithValue 传递数据 ===")
	ctx5 := context.WithValue(ctx, "request_id", "req-12345")
	ctx5 = context.WithValue(ctx5, "user_id", 42)

	processRequest(ctx5)

	// ============================================================
	//  6. context 链式传递
	// ============================================================
	fmt.Println("\n=== context 链式传递 ===")
	// 根 -> WithValue -> WithTimeout -> WithCancel
	root := context.Background()
	valCtx := context.WithValue(root, "trace_id", "abc-xyz")
	timeoutCtx, cancel := context.WithTimeout(valCtx, 3*time.Second)
	defer cancel()

	fmt.Printf("trace_id: %v\n", timeoutCtx.Value("trace_id"))
	dl, ok := timeoutCtx.Deadline()
	fmt.Printf(" deadline: %v (ok=%v)\n", dl, ok)
	fmt.Printf("    剩余: %v\n", time.Until(dl))

	fmt.Println("\n=== 总结 ===")
	fmt.Println("context.Background()  — 根 context")
	fmt.Println("WithCancel()          — 手动取消")
	fmt.Println("WithTimeout()         — 超时自动取消")
	fmt.Println("WithDeadline()        — 到截止时间取消")
	fmt.Println("WithValue()           — 传递请求数据")
	fmt.Println("ctx.Done()            — 监听取消信号")
	fmt.Println("ctx.Err()             — 获取取消原因")
}

func processRequest(ctx context.Context) {
	reqID := ctx.Value("request_id")
	userID := ctx.Value("user_id")
	fmt.Printf("处理请求: request_id=%v, user_id=%v\n", reqID, userID)
}

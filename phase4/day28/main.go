package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============================================================
//  并发下载器：
//  - 多个 worker 并发下载任务
//  - 支持超时控制（context）
//  - 实时进度显示
//  - 统计成功/失败/耗时
// ============================================================

type DownloadTask struct {
	ID   int
	URL  string
	Size int // MB
}

type DownloadResult struct {
	TaskID  int
	URL     string
	Success bool
	ErrMsg  string
	Duration time.Duration
}

// 模拟下载（本地演示，不真正发起网络请求）
func download(ctx context.Context, task DownloadTask) DownloadResult {
	start := time.Now()
	duration := time.Duration(task.Size*20) * time.Millisecond // 1MB = 20ms

	select {
	case <-time.After(duration):
		return DownloadResult{
			TaskID: task.ID, URL: task.URL, Success: true,
			Duration: time.Since(start),
		}
	case <-ctx.Done():
		return DownloadResult{
			TaskID: task.ID, URL: task.URL, Success: false,
			ErrMsg: ctx.Err().Error(), Duration: time.Since(start),
		}
	}
}

func main() {
	// 准备任务（模拟 8 个下载任务）
	tasks := make([]DownloadTask, 0)
	for i := 1; i <= 8; i++ {
		tasks = append(tasks, DownloadTask{
			ID:   i,
			URL:  fmt.Sprintf("https://example.com/file%d.zip", i),
			Size: 10 + rand.Intn(40), // 10-50 MB
		})
	}
	totalSize := 0
	for _, t := range tasks {
		totalSize += t.Size
	}
	fmt.Printf("准备下载 %d 个文件，总大小 %d MB\n", len(tasks), totalSize)

	// 全局超时：3 秒
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 启动 worker pool（3 个 worker）
	const numWorkers = 3
	taskChan := make(chan DownloadTask, len(tasks))
	resultChan := make(chan DownloadResult, len(tasks))

	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for task := range taskChan { // 从任务队列取任务，直到 channel 关闭
				fmt.Printf("[worker-%d] 开始下载 file%d.zip (%dMB)...\n",
					workerID, task.ID, task.Size)
				result := download(ctx, task)
				resultChan <- result
			}
		}(w)
	}

	// 发送所有任务，然后关闭 taskChan
	go func() {
		for _, t := range tasks {
			taskChan <- t
		}
		close(taskChan)
	}()

	// 等待所有 worker 完成后关闭 resultChan
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// 收集结果（带进度显示）
	start := time.Now()
	results := make([]DownloadResult, 0)
	done := 0
	downloadedMB := 0

	for result := range resultChan {
		results = append(results, result)
		done++

		if result.Success {
			downloadedMB += taskSize(tasks, result.TaskID)
			fmt.Printf("  ✓ file%d.zip 完成 (%v) [%d/%d, 已下载 %dMB]\n",
				result.TaskID, result.Duration.Round(time.Millisecond), done, len(tasks), downloadedMB)
		} else {
			fmt.Printf("  ✗ file%d.zip 失败: %s [%d/%d]\n",
				result.TaskID, result.ErrMsg, done, len(tasks))
		}
	}

	// 统计
	success := 0
	failed := 0
	var totalTime time.Duration
	for _, r := range results {
		if r.Success {
			success++
			totalTime += r.Duration
		} else {
			failed++
		}
	}

	fmt.Println("\n========== 下载统计 ==========")
	fmt.Printf("总耗时:   %v（串行预计需要 %v，加速比 %.1fx）\n",
		time.Since(start).Round(time.Millisecond),
		time.Duration(totalSize*20)*time.Millisecond,
		float64(time.Duration(totalSize*20)*time.Millisecond)/float64(time.Since(start)))
	fmt.Printf("成功:     %d 个\n", success)
	fmt.Printf("失败:     %d 个\n", failed)
	fmt.Printf("总大小:   %d MB / %d MB\n", downloadedMB, totalSize)
}

func taskSize(tasks []DownloadTask, id int) int {
	for _, t := range tasks {
		if t.ID == id {
			return t.Size
		}
	}
	return 0
}

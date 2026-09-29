package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// === HTTP 处理器 ===

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] %s %s\n", time.Now().Format("15:04:05"), r.Method, r.URL.Path)
	fmt.Fprintf(w, "Hello, %s!\n", r.URL.Query().Get("name"))
}

func jsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"message":"ok","method":"%s"}`+"\n", r.Method)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	// 获取路径参数（简单版）
	// /user/123 -> 提取 123
	fmt.Fprintf(w, "用户ID: %s\n", r.URL.Path[len("/user/"):])
}

func main() {
	// ============================================================
	//  1. 注册路由
	// ============================================================
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/api/status", jsonHandler)
	http.HandleFunc("/user/", userHandler)

	// ============================================================
	//  2. 启动服务器（在 goroutine 中，不阻塞）
	// ============================================================
	server := &http.Server{
		Addr:    ":8080",
		Handler: nil, // 使用默认的 DefaultServeMux
	}

	go func() {
		fmt.Println("🚀 HTTP 服务器启动在 http://localhost:8080")
		fmt.Println("可访问:")
		fmt.Println("  http://localhost:8080/?name=Go")
		fmt.Println("  http://localhost:8080/api/status")
		fmt.Println("  http://localhost:8080/user/123")
		fmt.Println("按 Ctrl+C 停止")
		server.ListenAndServe()
	}()

	// 等待服务器启动
	time.Sleep(200 * time.Millisecond)

	// ============================================================
	//  3. HTTP 客户端请求
	// ============================================================
	fmt.Println("\n=== HTTP 客户端请求 ===")

	// GET 请求
	resp, err := http.Get("http://localhost:8080/?name=Go学习者")
	if err != nil {
		fmt.Println("请求失败:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("GET /: %s", string(body))

	// POST 请求
	resp2, _ := http.Post("http://localhost:8080/api/status", "application/json", nil)
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	fmt.Printf("POST /api/status: %s", string(body2))

	// 自定义请求
	req, _ := http.NewRequest("GET", "http://localhost:8080/user/42", nil)
	client := &http.Client{Timeout: 5 * time.Second}
	resp3, _ := client.Do(req)
	defer resp3.Body.Close()
	body3, _ := io.ReadAll(resp3.Body)
	fmt.Printf("GET /user/42: %s", string(body3))

	// 保持运行一段时间让用户看到
	fmt.Println("\n服务器将在 3 秒后关闭...")
	time.Sleep(3 * time.Second)
	server.Close()
	fmt.Println("✓ 服务器已关闭")
}

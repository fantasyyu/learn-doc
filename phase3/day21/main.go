package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

// === 数据模型 ===
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// === 内存数据存储 ===
var users = []User{
	{ID: 1, Name: "Alice", Age: 25},
	{ID: 2, Name: "Bob", Age: 30},
}
var nextID = 3

// === 中间件：日志记录 ===
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next(w, r)
		log.Printf("[%s] %s %s %v", r.Method, r.URL.Path, r.RemoteAddr, time.Since(start))
	}
}

// === 中间件：CORS ===
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// === JSON 响应辅助函数 ===
func jsonResponse(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{Code: code, Message: message, Data: data})
}

// === API 处理器 ===

// GET /users — 获取所有用户
func getUsers(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, "success", users)
}

// GET /user?id=1 — 获取单个用户
func getUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "invalid id", nil)
		return
	}
	for _, u := range users {
		if u.ID == id {
			jsonResponse(w, http.StatusOK, "success", u)
			return
		}
	}
	jsonResponse(w, http.StatusNotFound, "user not found", nil)
}

// POST /user — 创建用户
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		jsonResponse(w, http.StatusBadRequest, "invalid json", nil)
		return
	}
	user.ID = nextID
	nextID++
	users = append(users, user)
	jsonResponse(w, http.StatusCreated, "created", user)
}

// DELETE /user?id=1 — 删除用户
func deleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, "invalid id", nil)
		return
	}
	for i, u := range users {
		if u.ID == id {
			users = append(users[:i], users[i+1:]...)
			jsonResponse(w, http.StatusOK, "deleted", nil)
			return
		}
	}
	jsonResponse(w, http.StatusNotFound, "user not found", nil)
}

// GET /health — 健康检查
func healthCheck(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, "ok", map[string]string{"status": "running"})
}

// === 路由分发 ===
func router(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/users":
		if r.Method == "GET" {
			getUsers(w, r)
		} else {
			jsonResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		}
	case "/user":
		switch r.Method {
		case "GET":
			getUser(w, r)
		case "POST":
			createUser(w, r)
		case "DELETE":
			deleteUser(w, r)
		default:
			jsonResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		}
	case "/health":
		healthCheck(w, r)
	default:
		jsonResponse(w, http.StatusNotFound, "not found", nil)
	}
}

func main() {
	// 命令行参数
	port := flag.Int("port", 8080, "服务器端口")
	flag.Parse()

	// 设置日志格式
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 注册路由（应用中间件）
	http.HandleFunc("/", corsMiddleware(loggingMiddleware(router)))

	// 创建服务器
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// 优雅关闭：监听系统信号
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		log.Println("→ 收到关闭信号，正在优雅关闭...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()

	// 启动服务器
	log.Printf("🚀 服务器启动在 http://localhost:%d", *port)
	log.Println("API 端点:")
	log.Println("  GET    /health     — 健康检查")
	log.Println("  GET    /users      — 用户列表")
	log.Println("  GET    /user?id=1  — 获取用户")
	log.Println("  POST   /user       — 创建用户")
	log.Println("  DELETE /user?id=1  — 删除用户")
	log.Println("按 Ctrl+C 停止")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("服务器错误: %v", err)
	}
	log.Println("✓ 服务器已关闭")
}

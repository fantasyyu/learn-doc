package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	baseURL := "http://localhost:8081"

	// 等待服务器启动
	time.Sleep(500 * time.Millisecond)

	fmt.Println("=== API 测试 ===")

	// 1. 健康检查
	resp, _ := http.Get(baseURL + "/health")
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("GET /health: %s\n", string(body))

	// 2. 获取用户列表
	resp, _ = http.Get(baseURL + "/users")
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("GET /users: %s\n", string(body))

	// 3. 创建用户
	newUser := map[string]interface{}{"name": "Charlie", "age": 28}
	jsonData, _ := json.Marshal(newUser)
	resp, _ = http.Post(baseURL+"/user", "application/json", bytes.NewBuffer(jsonData))
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("POST /user: %s\n", string(body))

	// 4. 获取单个用户
	resp, _ = http.Get(baseURL + "/user?id=1")
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("GET /user?id=1: %s\n", string(body))

	// 5. 删除用户
	req, _ := http.NewRequest("DELETE", baseURL+"/user?id=2", nil)
	resp, _ = http.DefaultClient.Do(req)
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("DELETE /user?id=2: %s\n", string(body))

	// 6. 再次获取用户列表
	resp, _ = http.Get(baseURL + "/users")
	body, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	fmt.Printf("GET /users (after delete): %s\n", string(body))

	fmt.Println("\n✅ 所有 API 测试完成！")
}

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// === 1. 基本结构体 ===
type User struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email,omitempty"` // omitempty: 空值时省略
	Age      int       `json:"age"`
	IsActive bool      `json:"is_active"`
	Created  time.Time `json:"created"`
	Password string    `json:"-"` // -: 完全忽略该字段
}

// === 2. 自定义序列化 ===
type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// MarshalJSON 自定义 JSON 输出格式
func (p Product) MarshalJSON() ([]byte, error) {
	type Alias Product // 避免递归
	return json.Marshal(&struct {
		Alias
		Price string `json:"price"` // 价格格式化为字符串
	}{
		Alias: (Alias)(p),
		Price: fmt.Sprintf("%.2f", p.Price),
	})
}

// === 3. 处理动态 JSON ===
func handleDynamicJSON() {
	fmt.Println("\n=== 动态 JSON（map[string]interface{}）===")

	jsonStr := `{"name":"Go","version":1.22,"features":["fast","simple"]}`

	// 解析到 map
	var result map[string]interface{}
	json.Unmarshal([]byte(jsonStr), &result)
	fmt.Printf("name: %v (type: %T)\n", result["name"], result["name"])
	fmt.Printf("version: %v (type: %T)\n", result["version"], result["version"])

	// 注意：数字默认解析为 float64！
	version := result["version"].(float64)
	fmt.Printf("version as float64: %.2f\n", version)
}

func main() {
	// ============================================================
	//  1. 结构体 -> JSON（序列化）
	// ============================================================
	fmt.Println("=== 结构体序列化 ===")

	user := User{
		ID:       1,
		Name:     "Alice",
		Email:    "alice@example.com",
		Age:      25,
		IsActive: true,
		Created:  time.Now(),
		Password: "secret123", // 会被忽略
	}

	data, _ := json.Marshal(user)
	fmt.Printf("Marshal:   %s\n", string(data))

	// 美化输出
	pretty, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("MarshalIndent:\n%s\n", string(pretty))

	// ============================================================
	//  2. JSON -> 结构体（反序列化）
	// ============================================================
	fmt.Println("\n=== JSON 反序列化 ===")

	jsonStr := `{"id":2,"name":"Bob","age":30,"is_active":false}`
	var u User
	err := json.Unmarshal([]byte(jsonStr), &u)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}
	fmt.Printf("解析结果: %+v\n", u)

	// ============================================================
	//  3. 切片序列化
	// ============================================================
	fmt.Println("\n=== 切片序列化 ===")

	users := []User{
		{ID: 1, Name: "Alice", Age: 25},
		{ID: 2, Name: "Bob", Age: 30},
	}
	usersJSON, _ := json.Marshal(users)
	fmt.Printf("用户列表: %s\n", string(usersJSON))

	// ============================================================
	//  4. omitempty 效果
	// ============================================================
	fmt.Println("\n=== omitempty 效果 ===")

	userNoEmail := User{ID: 3, Name: "Charlie", Age: 28}
	data2, _ := json.Marshal(userNoEmail)
	fmt.Printf("Email为空时省略: %s\n", string(data2))

	// ============================================================
	//  5. 自定义序列化
	// ============================================================
	fmt.Println("\n=== 自定义序列化 ===")

	p := Product{Name: "Book", Price: 29.99}
	pJSON, _ := json.Marshal(p)
	fmt.Printf("Product: %s\n", string(pJSON))

	// ============================================================
	//  6. 动态 JSON
	// ============================================================
	handleDynamicJSON()

	// ============================================================
	//  7. json.NewEncoder / json.NewDecoder（流式处理）
	// ============================================================
	fmt.Println("\n=== Encoder/Decoder 流式处理 ===")

	// Encoder 直接写到 io.Writer（如 http.ResponseWriter）
	encoder := json.NewEncoder(osStdout())
	encoder.Encode(map[string]string{"message": "hello"})
}

// osStdout 返回 os.Stdout 的接口（简化版）
func osStdout() interface{ Write([]byte) (int, error) } {
	return &stdoutWriter{}
}

type stdoutWriter struct{}

func (w *stdoutWriter) Write(p []byte) (int, error) {
	fmt.Print(string(p))
	return len(p), nil
}

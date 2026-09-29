package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	// ============================================================
	//  1. fmt 包 — 格式化输出
	// ============================================================
	fmt.Println("=== fmt 格式化 ===")

	name := "Go"
	version := 1.22
	isAwesome := true

	// Printf 格式化动词
	fmt.Printf("字符串: %s\n", name)
	fmt.Printf("整数: %d, 二进制: %b, 十六进制: %x\n", 255, 255, 255)
	fmt.Printf("浮点: %f, 科学计数: %e\n", version, version)
	fmt.Printf("布尔: %t\n", isAwesome)
	fmt.Printf("类型: %T\n", name)
	fmt.Printf("值: %v\n", name)   // 默认格式
	fmt.Printf("详细: %+v\n", struct{ A int }{A: 1}) // 显示字段名
	fmt.Printf("Go语法: %#v\n", name) // "Go"
	fmt.Printf("指针: %p\n", &name)

	// Sprintf 返回字符串
	s := fmt.Sprintf("Go %.2f", version)
	fmt.Printf("Sprintf 结果: %s\n", s)

	// ============================================================
	//  2. strings 包 — 字符串操作
	// ============================================================
	fmt.Println("\n=== strings 操作 ===")

	text := "  Hello, Go World!  "
	fmt.Printf("原始: %q\n", text)
	fmt.Printf("TrimSpace: %q\n", strings.TrimSpace(text))
	fmt.Printf("Contains 'Go': %v\n", strings.Contains(text, "Go"))
	fmt.Printf("HasPrefix 'Hello': %v\n", strings.HasPrefix(text, "Hello"))
	fmt.Printf("ToLower: %q\n", strings.ToLower(text))
	fmt.Printf("Replace: %q\n", strings.Replace(text, "Go", "Golang", 1))
	fmt.Printf("Split: %v\n", strings.Split("a,b,c", ","))
	fmt.Printf("Join: %q\n", strings.Join([]string{"a", "b", "c"}, "-"))

	// strings.Builder — 高效拼接字符串
	var b strings.Builder
	b.WriteString("Hello")
	b.WriteString(" ")
	b.WriteString("World")
	fmt.Printf("Builder: %q\n", b.String())

	// ============================================================
	//  3. strconv 包 — 字符串与数字转换
	// ============================================================
	fmt.Println("\n=== strconv 转换 ===")

	// 字符串转数字
	n, _ := strconv.Atoi("42")
	fmt.Printf("Atoi: %d (type: %T)\n", n, n)

	n2, _ := strconv.ParseInt("FF", 16, 64)
	fmt.Printf("ParseInt(FF, 16): %d\n", n2)

	f, _ := strconv.ParseFloat("3.14159", 64)
	fmt.Printf("ParseFloat: %.5f\n", f)

	// 数字转字符串
	s1 := strconv.Itoa(100)
	fmt.Printf("Itoa: %q\n", s1)

	s2 := strconv.FormatInt(255, 16)
	fmt.Printf("FormatInt(255, 16): %q\n", s2)

	// 布尔转字符串
	boolStr := strconv.FormatBool(true)
	fmt.Printf("FormatBool: %q\n", boolStr)

	// ============================================================
	//  4. time 包 — 时间处理
	// ============================================================
	fmt.Println("\n=== time 时间处理 ===")

	now := time.Now()
	fmt.Printf("当前时间: %v\n", now)
	fmt.Printf("格式化: %s\n", now.Format("2006-01-02 15:04:05"))
	fmt.Printf("时间戳(秒): %d\n", now.Unix())
	fmt.Printf("时间戳(毫秒): %d\n", now.UnixMilli())

	// 解析时间字符串
	parsed, _ := time.Parse("2006-01-02", "2024-09-04")
	fmt.Printf("解析: %v\n", parsed)

	// 时间运算
	tomorrow := now.Add(24 * time.Hour)
	fmt.Printf("明天: %s\n", tomorrow.Format("2006-01-02"))

	duration := tomorrow.Sub(now)
	fmt.Printf("间隔: %v\n", duration)

	// 定时器
	fmt.Println("\n=== 定时器 ===")
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	count := 0
	for t := range ticker.C {
		fmt.Printf("Tick %d: %v\n", count+1, t.Format("15:04:05.000"))
		count++
		if count >= 3 {
			break
		}
	}

	// time.Sleep
	fmt.Println("Sleep 500ms...")
	start := time.Now()
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("实际睡了: %v\n", time.Since(start))
}

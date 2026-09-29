package main

import (
	"fmt"
)

// init 函数：每个文件可以有多个，按文件名字母顺序执行
func init() {
	fmt.Println("init 1: 包初始化（最先执行）")
}

func init() {
	fmt.Println("init 2: 第二个 init")
}

// === 可见性规则 ===
// 大写开头 = public（包外可访问）
// 小写开头 = private（包内私有）

type PublicStruct struct { // 公开
	PublicField  string // 公开
	privateField string // 私有
}

func PublicFunc() {} // 公开
func privateFunc() {} // 私有

// === panic 和 recover ===
func safeDivide(a, b int) (result int, errMsg string) {
	// defer + recover 捕获 panic，防止程序崩溃
	defer func() {
		if r := recover(); r != nil {
			errMsg = fmt.Sprintf("捕获 panic: %v", r)
			result = 0
		}
	}()

	result = a / b // 如果 b=0，这里会 panic
	return result, ""
}

func riskyFunction() {
	panic("出大事了！") // 主动触发 panic
}

func safeCall() (errMsg string) {
	defer func() {
		if r := recover(); r != nil {
			errMsg = fmt.Sprintf("恢复成功: %v", r)
		}
	}()
	riskyFunction()
	return ""
}

func main() {
	fmt.Println("\n=== main 函数开始 ===")

	fmt.Println("\n--- panic / recover ---")
	// 正常情况
	r, err := safeDivide(10, 2)
	fmt.Printf("10/2 = %d, err=%q\n", r, err)

	// panic 被捕获
	r, err = safeDivide(10, 0)
	fmt.Printf("10/0 = %d, 消息=%q\n", r, err)

	// 主动 panic 被恢复
	msg := safeCall()
	fmt.Println(msg)

	fmt.Println("\n程序继续运行，没有崩溃！")

	fmt.Println("\n--- 可见性 ---")
	fmt.Println("大写开头 = public，小写开头 = private")
	fmt.Println("可见性在包级别，不是结构体级别！")

	fmt.Println("\n=== 使用建议 ===")
	fmt.Println("✅ 用 error 处理预期内的错误（文件不存在、网络超时）")
	fmt.Println("❌ 不要用 panic 控制正常流程")
	fmt.Println("✅ panic 只用于：不可恢复的错误（数组越界、空指针）")
	fmt.Println("✅ recover 一般只在顶层调用（如 HTTP handler）")
}

package main

import (
	"fmt"
	"time"
)

// calculate 四则运算
func calculate(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, fmt.Errorf("除数不能为0")
		}
		return a / b, nil
	default:
		return 0, fmt.Errorf("不支持的运算符: %s", op)
	}
}

func main() {
	start := time.Now()
	defer func() {
		fmt.Printf("\n运行耗时: %v\n", time.Since(start))
	}()

	fmt.Println("===== Go 命令行计算器（自动测试版）=====")

	// 测试用例: [a, b, op, 期望结果, 是否有错]
	tests := []struct {
		a, b      float64
		op        string
		expect    float64
		wantError bool
	}{
		{10, 5, "+", 15, false},
		{10, 5, "-", 5, false},
		{10, 5, "*", 50, false},
		{10, 5, "/", 2, false},
		{10, 0, "/", 0, true},   // 除0错误
		{10, 5, "^", 0, true},   // 不支持
	}

	history := make([]string, 0, len(tests))

	for _, tt := range tests {
		result, err := calculate(tt.a, tt.b, tt.op)
		expr := fmt.Sprintf("%.0f %s %.0f", tt.a, tt.op, tt.b)

		if tt.wantError {
			if err != nil {
				fmt.Printf("✓ %s → 错误: %v\n", expr, err)
			} else {
				fmt.Printf("✗ %s → 期望错误但得到 %.2f\n", expr, result)
			}
		} else {
			if err != nil {
				fmt.Printf("✗ %s → 意外错误: %v\n", expr, err)
			} else if result == tt.expect {
				record := fmt.Sprintf("%s = %.0f", expr, result)
				fmt.Printf("✓ %s\n", record)
				history = append(history, record)
			} else {
				fmt.Printf("✗ %s → 期望 %.0f, 得到 %.2f\n", expr, tt.expect, result)
			}
		}
	}

	fmt.Println("\n--- 成功记录 ---")
	for _, h := range history {
		fmt.Println(h)
	}

	fmt.Println("\n计算器完整代码: day07/calculator.go")
	fmt.Println("本地运行: go run day07/calculator.go")
}

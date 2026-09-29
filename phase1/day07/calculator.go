package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Calculator 实现四则运算，返回结果和错误
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
	// defer 记录程序运行时长
	start := time.Now()
	defer func() {
		fmt.Printf("\n程序运行耗时: %v\n", time.Since(start))
	}()

	fmt.Println("===== Go 命令行计算器 =====")
	fmt.Println("支持: + - * /    输入 'q' 退出")
	fmt.Println("示例: 10 + 5")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	history := make([]string, 0, 10) // 记录历史，切片

	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("读取输入失败:", err)
			continue
		}
		input = strings.TrimSpace(input)

		if input == "q" || input == "quit" {
			break
		}

		if input == "history" {
			fmt.Println("--- 历史记录 ---")
			for _, h := range history {
				fmt.Println(h)
			}
			continue
		}

		// 解析输入: "10 + 5"
		parts := strings.Fields(input)
		if len(parts) != 3 {
			fmt.Println("格式错误，请按: 数字 运算符 数字")
			continue
		}

		a, err1 := strconv.ParseFloat(parts[0], 64)
		b, err2 := strconv.ParseFloat(parts[2], 64)
		if err1 != nil || err2 != nil {
			fmt.Println("请输入有效的数字")
			continue
		}

		result, err := calculate(a, b, parts[1])
		if err != nil {
			fmt.Println("错误:", err)
			continue
		}

		record := fmt.Sprintf("%s = %.4f", input, result)
		fmt.Println("结果:", record)
		history = append(history, record)
	}

	fmt.Println("再见！")
}

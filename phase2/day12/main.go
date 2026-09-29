package main

import (
	"errors"
	"fmt"
)

// === 自定义错误类型 ===
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("字段 '%s' 验证失败: %s", e.Field, e.Message)
}

// === 带错误码的错误 ===
type BusinessError struct {
	Code    int
	Message string
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("[ERR-%d] %s", e.Code, e.Message)
}

// === 业务函数 ===
func divide(a, b float64) (float64, error) {
	if b == 0 {
		// errors.New 创建简单错误
		return 0, errors.New("除数不能为 0")
	}
	return a / b, nil
}

func validateAge(age int) error {
	if age < 0 {
		return &ValidationError{Field: "age", Message: "不能为负数"}
	}
	if age > 150 {
		return &ValidationError{Field: "age", Message: "超出人类范围"}
	}
	return nil
}

func withdraw(balance, amount float64) (float64, error) {
	if amount > balance {
		return balance, &BusinessError{Code: 1001, Message: "余额不足"}
	}
	return balance - amount, nil
}

// === 错误包装（Go 1.13+）===
func process() error {
	if _, err := divide(10, 0); err != nil {
		// fmt.Errorf + %w 包装错误，保留原始错误链
		return fmt.Errorf("处理失败: %w", err)
	}
	return nil
}

func main() {
	fmt.Println("=== 基础错误处理 ===")
	result, err := divide(10, 0)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Println("结果:", result)
	}

	fmt.Println("\n=== 自定义错误类型 ===")
	if err := validateAge(-5); err != nil {
		fmt.Println(err)
		// 类型断言判断具体错误类型
		if ve, ok := err.(*ValidationError); ok {
			fmt.Printf("是 ValidationError，字段: %s\n", ve.Field)
		}
	}

	fmt.Println("\n=== 业务错误 ===")
	_, err = withdraw(100, 200)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("\n=== 错误包装 ===")
	err = process()
	if err != nil {
		fmt.Println("包装后:", err)
		// errors.Unwrap 解包
		fmt.Println("原始错误:", errors.Unwrap(err))
	}

	fmt.Println("\n=== errors.Is 判断错误链 ===")
	baseErr := errors.New("连接超时")
	wrappedErr := fmt.Errorf("服务A调用失败: %w", baseErr)
	fmt.Println("errors.Is(wrappedErr, baseErr) =", errors.Is(wrappedErr, baseErr))
}

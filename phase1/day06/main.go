package main

import (
	"errors"
	"fmt"
)

// === 1. 多返回值（Go 标志性特性）===
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("除数不能为0")
	}
	return a / b, nil
}

// === 2. 命名返回值 ===
func rectangle(w, h float64) (area, perimeter float64) {
	area = w * h         // 直接赋值给命名返回值
	perimeter = 2 * (w + h)
	return // 裸 return，自动返回命名变量
}

// === 3. 变参函数 ===
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// === 4. defer 后进先出（LIFO）===
func deferDemo() {
	fmt.Println("\n--- defer 执行顺序 ---")
	defer fmt.Println("defer 1")
	defer fmt.Println("defer 2")
	defer fmt.Println("defer 3")
	fmt.Println("函数体执行完毕")
}

// defer 常用于：关闭文件、释放锁、记录执行时间
func deferWithParam() {
	fmt.Println("\n--- defer 参数求值时机 ---")
	x := 1
	defer fmt.Printf("defer x=%d\n", x) // 此时 x=1 已被捕获
	x = 100
	fmt.Printf("函数内 x=%d\n", x)
}

func main() {
	// 多返回值
	result, err := divide(10, 3)
	if err != nil {
		fmt.Println("错误:", err)
	} else {
		fmt.Printf("10/3 = %.2f\n", result)
	}

	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("错误:", err)
	}

	// 命名返回值
	a, p := rectangle(3, 4)
	fmt.Printf("矩形: 面积=%.0f, 周长=%.0f\n", a, p)

	// 变参
	fmt.Printf("sum(1,2,3,4) = %d\n", sum(1, 2, 3, 4))
	fmt.Printf("sum() = %d\n", sum())
	// 切片展开传入
	nums := []int{10, 20, 30}
	fmt.Printf("sum(slice...) = %d\n", sum(nums...))

	// defer
	deferDemo()
	deferWithParam()

	// === 5. 函数作为值（一等公民）===
	fmt.Println("\n--- 函数作为参数 ---")
	operate := func(a, b int, op func(int, int) int) int {
		return op(a, b)
	}
	add := func(a, b int) int { return a + b }
	mul := func(a, b int) int { return a * b }
	fmt.Printf("operate(3,4,add)=%d\n", operate(3, 4, add))
	fmt.Printf("operate(3,4,mul)=%d\n", operate(3, 4, mul))
}

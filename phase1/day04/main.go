package main

import "fmt"

func main() {
	// === for 的三种形态 ===

	// 1. 经典 C-style for
	for i := 0; i < 3; i++ {
		fmt.Println("循环i:", i)
	}

	// 2. 条件 for（相当于 while）
	j := 0
	for j < 3 {
		fmt.Println("循环j:", j)
		j++
	}

	// 3. 无限循环（相当于 for(;;) / while true）
	k := 0
	for {
		if k >= 2 {
			break
		}
		fmt.Println("无限循环k:", k)
		k++
	}

	// === if 支持初始化语句 ===
	if x := 10; x > 5 {
		fmt.Println("x > 5, x =", x)
	} // x 的作用域仅限 if 块

	// === switch ===
	grade := "B"
	switch grade {
	case "A":
		fmt.Println("优秀")
	case "B":
		fmt.Println("良好") // 匹配到这里
	case "C":
		fmt.Println("及格")
	default:
		fmt.Println("不及格")
	}

	// switch 可替代多 if-else，且默认 break（不需要写）
	score := 85
	switch {
	case score >= 90:
		fmt.Println("A")
	case score >= 80:
		fmt.Println("B") // 匹配到这里
	case score >= 60:
		fmt.Println("C")
	}

	// === for + range（后面 day5 重点讲）===
	str := "Go"
	for i, ch := range str {
		fmt.Printf("index=%d char=%c code=%d\n", i, ch, ch)
	}
}

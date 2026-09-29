package main

import "fmt"

func main() {
	// === 1. 变量声明 ===
	var a int = 10       // 完整声明
	var b = 20           // 类型推断
	c := 30              // 短变量声明（只能在函数内！）
	fmt.Printf("a=%d b=%d c=%d\n", a, b, c)

	// === 2. 多个变量 ===
	var x, y = 1, 2
	m, n := "hello", 3.14
	fmt.Printf("x=%d y=%d m=%s n=%.2f\n", x, y, m, n)

	// === 3. 常量 ===
	const Pi = 3.14159
	const (
		Monday = iota // 0
		Tuesday       // 1
		Wednesday     // 2
	)
	fmt.Printf("Pi=%.5f, Tuesday=%d\n", Pi, Tuesday)

	// === 4. 基本类型速览 ===
	var i8 int8 = 127          // -128 ~ 127
	var ui uint = 0            // 无符号整型
	var f64 float64 = 1.234    // 默认浮点类型
	var s string = "Go"        // 字符串不可变
	var ok bool = true         // true/false
	var r rune = '中'          // rune = int32，表示 Unicode 码点
	var by byte = 'A'          // byte = uint8

	fmt.Printf("类型: i8=%T f64=%T s=%T ok=%v r=%T(%d) by=%T(%d) ui=%d\n",
		i8, f64, s, ok, r, r, by, by, ui)

	// === 5. 零值（未初始化时的默认值）===
	var zeroInt int
	var zeroStr string
	var zeroBool bool
	fmt.Printf("零值: int=%d string=%q bool=%v\n", zeroInt, zeroStr, zeroBool)
}

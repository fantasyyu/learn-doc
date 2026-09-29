package main

import "fmt"

// 形状接口
type Shape interface {
	Area() float64
}

type Rectangle struct {
	Width, Height float64
}

func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func main() {
	fmt.Println("=== 类型断言 ===")
	var s Shape = Rectangle{Width: 3, Height: 4}

	// 方式1：断言成功返回具体值 + true
	if rect, ok := s.(Rectangle); ok {
		fmt.Printf("是 Rectangle: %+v, 面积=%.0f\n", rect, rect.Area())
	}

	// 方式2：直接断言（失败会 panic！）
	// circ := s.(Circle) // ❌ 运行时 panic！

	fmt.Println("\n=== 类型开关 type switch ===")
	shapes := []Shape{
		Rectangle{2, 3},
		Circle{5},
		Rectangle{4, 5},
	}

	for _, shape := range shapes {
		switch v := shape.(type) {
		case Rectangle:
			fmt.Printf("矩形: 宽=%.0f 高=%.0f, 面积=%.0f\n", v.Width, v.Height, v.Area())
		case Circle:
			fmt.Printf("圆形: 半径=%.0f, 面积=%.2f\n", v.Radius, v.Area())
		default:
			fmt.Println("未知形状")
		}
	}

	fmt.Println("\n=== interface{} 转具体类型 ===")
	var anything interface{} = "hello"

	if str, ok := anything.(string); ok {
		fmt.Printf("字符串: %s, 长度: %d\n", str, len(str))
	}

	anything = 42
	if num, ok := anything.(int); ok {
		fmt.Printf("整数: %d, 平方: %d\n", num, num*num)
	}

	fmt.Println("\n=== 判断接口是否实现了其他接口 ===")
	var s2 Shape = Circle{1}
	if _, ok := s2.(fmt.Stringer); ok {
		fmt.Println("Circle 实现了 fmt.Stringer")
	} else {
		fmt.Println("Circle 没有实现 fmt.Stringer")
	}
}

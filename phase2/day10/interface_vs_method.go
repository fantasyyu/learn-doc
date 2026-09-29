package main

import "fmt"

// ============================================================
//  你的困惑：方法定义 和 接口实现 为什么长得一样？
// ============================================================

// === 1. 先定义接口 ===
type Speaker1 interface {
	Speak() string
}

type Mover1 interface {
	Move()
}

// === 2. 定义类型 + 绑定方法 ===
type Dog1 struct {
	Name string
}

// 这行代码是在做什么？
//
//	A. 定义 Dog 的方法？  -- 是的！
//	B. 实现 Speaker 接口？ -- 也是的！
//
// 但代码层面，它只做 A，B 是编译器事后推导出来的
func (d Dog1) Speak() string {
	return d.Name + "说汪汪"
}

func (d Dog1) Move() {
	fmt.Println(d.Name, "在跑")
}

// === 3. 另一个类型：只有 Speak，没有 Move ===
type Parrot struct {
	Name string
}

func (p Parrot) Speak() string {
	return p.Name + "说你好"
}

// Parrot 没有 Move() 方法，所以不满足 Mover 接口

// ============================================================
//  关键点：Go 的接口实现是 事后检查，不是 事先声明
// ============================================================

func test1() {
	fmt.Println("=== Go 的隐式实现 ===")

	dog := Dog{Name: "旺财"}
	parrot := Parrot{Name: "鹦鹉"}

	// 编译器在这里检查：
	// "Dog 有 Speak() 吗？有！-> 可以赋值给 Speaker"
	var s1 Speaker = dog
	fmt.Println(s1.Speak())

	// "Parrot 有 Speak() 吗？有！-> 也可以赋值给 Speaker"
	var s2 Speaker = parrot
	fmt.Println(s2.Speak())

	// "Dog 有 Move() 吗？有！-> 满足组合接口"
	var a interface {
		Speak() string
		Move()
	} = dog
	a.Move()

	fmt.Println("\n=== 编译时断言：显式声明我实现了某个接口 ===")
	// 虽然 Go 是隐式实现，但你可以用这种方式显式声明
	// 如果 Dog 没有实现 Speaker，这行代码会编译报错
	var _ Speaker = Dog{}
	var _ Mover = Dog{}
	fmt.Println("Dog 显式声明实现了 Speaker 和 Mover（编译期检查）")

	fmt.Println("\n=== 与 Java/C# 的对比 ===")
	fmt.Println("Java:  class Dog implements Speaker { ... }")
	fmt.Println("Go:    func (d Dog) Speak() string { ... }  <- 没有 implements！")
	fmt.Println("       编译器自己发现 Dog 有 Speak，于是认为它实现了 Speaker")

	fmt.Println("\n=== 这就是 Go 的 鸭子类型 ===")
	fmt.Println("如果它走起路来像鸭子，叫起来像鸭子，那它就是鸭子")
	fmt.Println("不需要显式说 我是一只鸭子")

	fmt.Println("\n=== 方法定义的语法，本质是在给类型附加行为 ===")
	fmt.Println("给 Counter 加一个 Show 方法，跟接口没关系")
	fmt.Println("除非后面有人定义了需要 Show 的接口，Counter 就自动满足了")
}

type Counter struct{ Value int }

func (c Counter) Show() {
	fmt.Printf("Counter.Value = %d\n", c.Value)
}

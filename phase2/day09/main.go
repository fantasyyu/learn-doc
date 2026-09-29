package main

import "fmt"

// Counter 计数器
type Counter struct {
	Value int
}

// === 值接收者：方法接收的是结构体的拷贝 ===
func (c Counter) Show() {
	fmt.Printf("[值接收者] Value = %d\n", c.Value)
}

func (c Counter) Increment() {
	c.Value++ // 只改了拷贝，原结构体不变！
	fmt.Printf("[值接收者] 内部 Value = %d\n", c.Value)
}

// === 指针接收者：方法接收的是结构体的地址 ===
func (c *Counter) ShowPtr() {
	fmt.Printf("[指针接收者] Value = %d\n", c.Value)
}

func (c *Counter) IncrementPtr() {
	c.Value++ // 修改的是原结构体！
	fmt.Printf("[指针接收者] 内部 Value = %d\n", c.Value)
}

// 指针接收者 + 修改大数据结构时更省内存（不拷贝整个结构体）
type BigStruct struct {
	Data [1000000]int
}

func (b BigStruct) ValueMethod() {
	fmt.Println("值接收者：传入了整个 BigStruct 的拷贝！")
}

func (b *BigStruct) PtrMethod() {
	fmt.Println("指针接收者：只传入了 8 字节的地址！")
}

func main() {
	fmt.Println("=== 值接收者 vs 指针接收者 ===\n")

	c := Counter{Value: 10}
	fmt.Printf("初始: c.Value = %d\n", c.Value)

	c.Increment()        // 值接收者
	fmt.Printf("调用后: c.Value = %d ← 没变！\n\n", c.Value)

	c.IncrementPtr()     // 指针接收者
	fmt.Printf("调用后: c.Value = %d ← 变了！\n\n", c.Value)

	fmt.Println("=== 自动解引用 ===")
	// Go 会自动处理：值变量可调指针方法，指针变量可调值方法
	c.Show()             // c 是值，但 Show() 是值接收者 → OK
	c.ShowPtr()          // c 是值，但 ShowPtr() 是指针接收者 → Go 自动取地址 &c

	cp := &Counter{Value: 20}
	cp.Show()            // cp 是指针，但 Show() 是值接收者 → Go 自动解引用 *cp
	cp.ShowPtr()         // cp 是指针，ShowPtr() 是指针接收者 → OK

	fmt.Println("\n=== 选择建议 ===")
	fmt.Println("1. 需要修改接收者 → 用指针接收者")
	fmt.Println("2. 结构体很大 → 用指针接收者（避免拷贝开销）")
	fmt.Println("3. 保证一致性 → 某个方法用指针，其他方法也用指针")
	fmt.Println("4. 不需要修改且结构体小 → 值接收者（线程安全）")
}

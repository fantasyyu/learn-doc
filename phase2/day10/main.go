package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ============================================================
//  接口的本质：
//  接口在底层是一个 "iface" 结构体：
//
//  type iface struct {
//      tab  *itab          // 接口表：包含类型信息 + 方法地址表
//      data unsafe.Pointer // 指向具体值的指针
//  }
//
//  也就是说，一个接口变量 = (类型信息, 值的指针)
//  这正是 "动态类型" 和 "动态值" 的来源
// ============================================================

// === 1. 定义接口 ===
type Speaker interface {
	Speak() string
}

type Mover interface {
	Move()
}

// === 2. 隐式实现：Dog 和 Cat 都自动满足 Speaker 接口 ===
type Dog struct {
	Name string
}

func (d Dog) Speak() string { // 值接收者
	return d.Name + "说：汪汪！"
}

func (d Dog) Move() {
	fmt.Println(d.Name, "在跑")
}

type Cat struct {
	Name string
}

func (c Cat) Speak() string { // 值接收者
	return c.Name + "说：喵喵~"
}

func (c Cat) Move() {
	fmt.Println(c.Name, "在跳")
}

// === 3. 多态：不关心具体类型，只关心接口 ===
func MakeSpeak(s Speaker) {
	fmt.Println(s.Speak())
}

// === 4. 空接口：interface{} 可以装任何类型 ===
func Describe(v interface{}) {
	fmt.Printf("Describe: 类型=%T, 值=%v\n", v, v)
}

// === 5. 接口组合 ===
type Animal interface {
	Speaker
	Mover
}

func ShowAnimal(a Animal) {
	a.Move()
	fmt.Println(a.Speak())
}

// ============================================================
//  ⭐ 面试重点 1：nil 接口的陷阱
// ============================================================

// MyError 自定义错误类型
type MyError struct {
	Msg string
}

func (e *MyError) Error() string {
	return e.Msg
}

// 这个函数有经典的 nil 接口陷阱！
func returnNilInterface() error {
	var err *MyError = nil // 具体类型的 nil 指针
	return err             // 返回后 err 不是 nil！
}

// ============================================================
//  ⭐ 面试重点 2：值接收者 vs 指针接收者 对接口的影响
// ============================================================

type Printer interface {
	Print()
}

type Document struct {
	Content string
}

// 值接收者
func (d Document) Print() {
	fmt.Println("Document:", d.Content)
}

type Report struct {
	Title string
}

// 指针接收者
func (r *Report) Print() {
	fmt.Println("Report:", r.Title)
}

// ============================================================
//  ⭐ 面试重点 3：接口在实际开发中的作用
// ============================================================

// Storage 接口 — 定义存储行为
type Storage interface {
	Save(key string, data string) error
	Load(key string) (string, error)
}

// MemoryStorage 内存实现（测试用）
type MemoryStorage struct {
	data map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{data: make(map[string]string)}
}

func (m *MemoryStorage) Save(key, data string) error {
	m.data[key] = data
	return nil
}

func (m *MemoryStorage) Load(key string) (string, error) {
	v, ok := m.data[key]
	if !ok {
		return "", fmt.Errorf("key %s not found", key)
	}
	return v, nil
}

// FileStorage 文件实现（生产用）
type FileStorage struct {
	Path string
}

func (f *FileStorage) Save(key, data string) error {
	fmt.Printf("[FileStorage] Save %s = %s\n", key, data)
	return nil
}

func (f *FileStorage) Load(key string) (string, error) {
	fmt.Printf("[FileStorage] Load %s\n", key)
	return "file-data", nil
}

// UserService 依赖 Storage 接口，不关心具体实现
type UserService struct {
	store Storage // 依赖接口，不是具体类型
}

func (us *UserService) CreateUser(id, name string) error {
	return us.store.Save(id, name)
}

func (us *UserService) GetUser(id string) (string, error) {
	return us.store.Load(id)
}

// ============================================================
//  常用标准库接口
// ============================================================

func demoStandardInterfaces() {
	fmt.Println("\n=== 标准库接口：io.Reader ===")
	// strings.Reader 实现了 io.Reader
	r := strings.NewReader("Hello, Go!")
	buf := make([]byte, 5)
	n, _ := r.Read(buf)
	fmt.Printf("Read %d bytes: %s\n", n, string(buf))

	fmt.Println("\n=== 标准库接口：io.Writer ===")
	// os.Stdout 实现了 io.Writer
	io.WriteString(os.Stdout, "写入到标准输出\n")
}

func main() {
	fmt.Println("=== 1. 隐式实现与多态 ===")
	dog := Dog{"旺财"}
	cat := Cat{"咪咪"}
	MakeSpeak(dog)
	MakeSpeak(cat)

	fmt.Println("\n=== 2. 接口组合 ===")
	ShowAnimal(dog)

	fmt.Println("\n=== 3. 空接口 interface{} ===")
	Describe(42)
	Describe("hello")
	Describe(dog)

	// ============================================================
	// ⭐ nil 接口陷阱详解
	// ============================================================
	fmt.Println("\n=== 4. ⭐ nil 接口陷阱 ===")

	// 情况 A：error 接口变量直接赋 nil -> 确实是 nil
	var err1 error = nil
	fmt.Printf("err1 == nil ? %v\n", err1 == nil) // true

	// 情况 B：具体类型的 nil 指针赋值给接口 -> 不是 nil！
	var err2 error = returnNilInterface()
	fmt.Printf("err2 == nil ? %v <- 坑！具体类型的 nil 指针传给接口后，接口不是 nil\n", err2 == nil)

	// 原因：err2 = (类型=*MyError, 数据指针=nil)
	// 接口变量本身是 (tab, data) 的组合，tab 不为 nil，所以接口变量 != nil

	// 正确写法：返回时直接写 nil
	funcCorrect := func() error {
		if false {
			return &MyError{"err"}
		}
		return nil // 直接返回 nil
	}
	fmt.Printf("正确写法: err == nil ? %v\n", funcCorrect() == nil)

	// ============================================================
	// ⭐ 值接收者 vs 指针接收者 对接口的影响
	// ============================================================
	fmt.Println("\n=== 5. ⭐ 值/指针接收者 对接口的影响 ===")

	doc := Document{"Hello"}
	var p1 Printer = doc    // 值接收者：值和指针都能满足接口
	var p2 Printer = &doc   // 值接收者：指针也能满足（自动解引用）
	p1.Print()
	p2.Print()

	report := Report{"Q3 Report"}
	// var p3 Printer = report  // 编译错误！指针接收者，值不能满足
	var p4 Printer = &report // 指针接收者：只有指针能满足
	p4.Print()

	fmt.Println("\n总结:")
	fmt.Println("  值接收者方法: 值变量和指针变量都能满足接口")
	fmt.Println("  指针接收者方法: 只有指针变量能满足接口")

	// ============================================================
	// 接口在实战中的应用：依赖注入
	// ============================================================
	fmt.Println("\n=== 6. 接口实战：依赖注入 ===")

	// 测试环境：用内存存储
	memStore := NewMemoryStorage()
	service := &UserService{store: memStore}
	service.CreateUser("001", "Alice")
	name, _ := service.GetUser("001")
	fmt.Printf("从内存读取: %s\n", name)

	// 生产环境：换成文件存储，UserService 代码完全不用改
	fileStore := &FileStorage{Path: "/data"}
	service2 := &UserService{store: fileStore}
	service2.CreateUser("002", "Bob")

	// ============================================================
	// 标准库接口
	// ============================================================
	demoStandardInterfaces()

	// ============================================================
	// 接口值的可比较性
	// ============================================================
	fmt.Println("\n=== 7. 接口值的可比较性 ===")
	var a Speaker = Dog{"A"}
	var b Speaker = Dog{"A"}
	var c Speaker = Cat{"A"}
	_ = c // 仅用于展示：不同动态类型的接口不能比较
	fmt.Printf("Dog{A} == Dog{A} ? %v\n", a == b) // true（动态类型和值都相同）
	// fmt.Println(a == c) // 编译错误！不同动态类型不能比较

	fmt.Println("\n=== 总结 ===")
	fmt.Println("1. 接口 = (动态类型, 动态值)")
	fmt.Println("2. nil 接口 != 具体类型的 nil 指针赋值给接口")
	fmt.Println("3. 值接收者方法: 值/指针都能满足接口")
	fmt.Println("4. 指针接收者方法: 只有指针能满足接口")
	fmt.Println("5. 接口用于解耦：依赖接口，而非具体实现")
}

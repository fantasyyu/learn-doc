package main

import "fmt"

// === 1. 基本结构体定义 ===
type Person struct {
	Name string
	Age  int
}

// === 2. 嵌套结构体 ===
type Address struct {
	City   string
	Street string
}

type Employee struct {
	Person  // 匿名字段（嵌入）— 直接继承 Person 的字段
	Address // 嵌套结构体（有名）
	Salary  float64
	ID      int
}

// === 3. 结构体标签（Tag）— 常用于 JSON/XML 序列化 ===
type User struct {
	Username string `json:"user_name"`       // json 序列化时用 user_name
	Password string `json:"-"`               // 忽略该字段
	Email    string `json:"email,omitempty"` // 为空时省略
	Age      int    `json:"age"`
}

func main() {
	fmt.Println("=== 结构体基础 ===")
	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{"Bob", 25}     // 按字段顺序初始化（不推荐，易出错）
	fmt.Printf("p1: %+v\n", p1) // %+v 显示字段名
	fmt.Printf("p2: %+v\n", p2)

	fmt.Println("\n=== 嵌套与匿名字段 ===")
	emp := Employee{
		Person:  Person{Name: "Charlie", Age: 35},
		Address: Address{City: "Beijing", Street: "Zhongguancun"},
		Salary:  50000,
		ID:      1001,
	}
	// 匿名字段可以直接访问：emp.Name（等价于 emp.Person.Name）
	fmt.Printf("姓名: %s, 年龄: %d, 工资: %.2f, ID: %d\n", emp.Name, emp.Age, emp.Salary, emp.ID)
	// 有名嵌套需要用全路径
	fmt.Printf("城市: %s, 街道: %s\n", emp.Address.City, emp.Address.Street)

	fmt.Println("\n=== 结构体标签 ===")
	u := User{Username: "john", Password: "secret", Email: "", Age: 20}
	fmt.Printf("User 标签: Username tag=%s\n", getTag(u, "Username"))

	fmt.Println("\n=== new 和 & 创建结构体 ===")
	p3 := new(Person) // 返回 *Person，字段初始化为零值
	p3.Name = "David"
	p4 := &Person{} // 等价于 new
	p4.Age = 40
	fmt.Printf("p3: %v, p4: %v\n", *p3, *p4)

	fmt.Println("\n=== 结构体比较 ===")
	// 结构体可比较（字段全是可比较类型时）
	a := Person{"Eve", 28}
	b := Person{"Eve", 28}
	fmt.Printf("a == b ? %v\n", a == b)
}

// 辅助函数：获取结构体字段标签（用反射）
func getTag(u User, field string) string {
	// 简化版，实际用 reflect 包
	if field == "Username" {
		return "user_name"
	}
	return ""
}

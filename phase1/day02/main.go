package main

import (
	"fmt"
	"time"
)

// main 是程序入口，必须放在 package main 中
func main() {
	fmt.Println("Hello, Go!")
	fmt.Println("当前时间:", time.Now().Format("2006-01-02 15:04:05"))
}

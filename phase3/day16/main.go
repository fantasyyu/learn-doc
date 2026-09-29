package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	// ============================================================
	//  1. os 包 — 文件基本操作
	// ============================================================
	fmt.Println("=== os 文件操作 ===")

	// 创建/写入文件
	content := []byte("Hello, Go!\n这是第二行。\n")
	err := os.WriteFile("test.txt", content, 0644)
	if err != nil {
		fmt.Println("写入失败:", err)
		return
	}
	fmt.Println("✓ 写入文件 test.txt")

	// 读取整个文件
	data, err := os.ReadFile("test.txt")
	if err != nil {
		fmt.Println("读取失败:", err)
		return
	}
	fmt.Printf("读取内容: %q\n", string(data))

	// 文件信息
	info, _ := os.Stat("test.txt")
	fmt.Printf("文件名: %s, 大小: %d bytes, 修改时间: %v\n",
		info.Name(), info.Size(), info.ModTime().Format("2006-01-02 15:04"))

	// ============================================================
	//  2. os 打开文件（低级别操作）
	// ============================================================
	fmt.Println("\n=== os.Open / os.Create ===")

	file, err := os.Open("test.txt")
	if err != nil {
		fmt.Println("打开失败:", err)
		return
	}
	defer file.Close() // 记得关闭！

	// 读取部分内容
	buf := make([]byte, 10)
	n, _ := file.Read(buf)
	fmt.Printf("Read %d bytes: %q\n", n, string(buf[:n]))

	// 追加写入
	f2, err := os.OpenFile("test.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("打开失败:", err)
		return
	}
	defer f2.Close()
	f2.WriteString("追加的一行\n")
	fmt.Println("✓ 追加写入")

	// ============================================================
	//  3. bufio — 带缓冲的读写（推荐用于大文件）
	// ============================================================
	fmt.Println("\n=== bufio 缓冲读写 ===")

	// 逐行读取（最常用的文件读取方式）
	f3, _ := os.Open("test.txt")
	defer f3.Close()

	scanner := bufio.NewScanner(f3)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		fmt.Printf("  行%d: %s\n", lineNum, scanner.Text())
	}

	// 带缓冲的写入
	f4, _ := os.Create("buffered.txt")
	defer f4.Close()

	writer := bufio.NewWriter(f4)
	writer.WriteString("第一行\n")
	writer.WriteString("第二行\n")
	writer.Flush() // 必须 Flush，否则数据还在缓冲区！
	fmt.Println("✓ 缓冲写入完成")

	// ============================================================
	//  4. io 包 — 通用读写接口
	// ============================================================
	fmt.Println("\n=== io 包通用接口 ===")

	// io.Copy：从一个 Reader 拷贝到 Writer
	src, _ := os.Open("test.txt")
	defer src.Close()
	dst, _ := os.Create("copy.txt")
	defer dst.Close()

	copied, _ := io.Copy(dst, src)
	fmt.Printf("✓ 拷贝 %d bytes\n", copied)

	// io.ReadAll：读取全部内容
	f5, _ := os.Open("test.txt")
	defer f5.Close()
	all, _ := io.ReadAll(f5)
	fmt.Printf("ReadAll: %d bytes\n", len(all))

	// ============================================================
	//  5. 清理
	// ============================================================
	fmt.Println("\n=== 清理临时文件 ===")
	os.Remove("buffered.txt")
	os.Remove("copy.txt")
	// test.txt 保留作为示例
	fmt.Println("✓ 清理完成，test.txt 保留")
}

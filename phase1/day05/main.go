package main

import "fmt"

func main() {
	// ========== 1. 数组（固定长度，值类型）==========
	var arr1 [3]int = [3]int{1, 2, 3}
	arr2 := [3]int{4, 5, 6}
	arr3 := [...]int{7, 8, 9, 10} // 编译器自动数长度

	fmt.Println("数组:", arr1, arr2, arr3)
	fmt.Printf("arr1 类型: %T, 长度: %d\n", arr1, len(arr1))

	// 数组是值类型，赋值会拷贝整个数组！
	arrCopy := arr1
	arrCopy[0] = 999
	fmt.Printf("arr1[0]=%d, arrCopy[0]=%d （值拷贝）\n", arr1[0], arrCopy[0])

	// ========== 2. 切片（Slice）— 重点！==========
	// 切片是引用类型，底层指向一个数组，包含：指针 + 长度(len) + 容量(cap)

	// 方式1：从数组切片
	arr := [5]int{10, 20, 30, 40, 50}
	s1 := arr[1:4] // 左闭右开，取 index 1,2,3
	fmt.Printf("\ns1=%v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
	// cap = 从切片的起始位置到数组末尾的元素数 = 5-1 = 4

	// 方式2：make 创建
	s2 := make([]int, 3, 5) // 长度3，容量5
	fmt.Printf("s2=%v, len=%d, cap=%d\n", s2, len(s2), cap(s2))

	// 方式3：字面量
	s3 := []int{100, 200, 300}
	fmt.Printf("s3=%v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// ⚠️ 关键：切片是引用类型，修改会影响底层数组
	s1[0] = 999
	fmt.Printf("修改 s1[0] 后，arr=%v （原数组变了！）\n", arr)

	// append — 长度超过 cap 时会触发扩容（分配新底层数组）
	s4 := []int{1, 2}
	fmt.Printf("\nappend 前: s4=%v, len=%d, cap=%d, ptr=%p\n", s4, len(s4), cap(s4), s4)
	s4 = append(s4, 3, 4)
	fmt.Printf("append 后: s4=%v, len=%d, cap=%d, ptr=%p\n", s4, len(s4), cap(s4), s4)

	// 扩容后 cap 不足，底层数组会变！
	s5 := make([]int, 2, 2)
	s5[0], s5[1] = 1, 2
	oldPtr := s5
	s5 = append(s5, 3) // 超过 cap，需要扩容
	fmt.Printf("\n扩容前后指针相同? %v\n", &oldPtr[0] == &s5[0])

	// copy 切片
	src := []int{1, 2, 3, 4, 5}
	dst := make([]int, 3)
	n := copy(dst, src)
	fmt.Printf("\ncopy: dst=%v, 复制了%d个元素\n", dst, n)

	// ========== 3. map（引用类型，无序）==========
	m := make(map[string]int)
	m["alice"] = 25
	m["bob"] = 30
	fmt.Printf("\nmap: %v\n", m)

	// 取值 + 判断是否存在
	if age, ok := m["alice"]; ok {
		fmt.Printf("alice 的年龄: %d\n", age)
	}
	if _, ok := m["charlie"]; !ok {
		fmt.Println("charlie 不存在")
	}

	// 删除
	delete(m, "bob")
	fmt.Printf("删除 bob 后: %v\n", m)

	// map 字面量
	scores := map[string]int{
		"数学": 90,
		"英语": 85,
	}
	fmt.Printf("scores: %v\n", scores)

	// ========== 4. range 遍历 ==========
	fmt.Println("\n--- range 遍历数组 ---")
	for i, v := range arr1 {
		fmt.Printf("index=%d, value=%d\n", i, v)
	}

	fmt.Println("\n--- range 遍历切片 ---")
	for i, v := range s3 {
		fmt.Printf("index=%d, value=%d\n", i, v)
	}

	fmt.Println("\n--- range 遍历 map（无序！）---")
	for k, v := range scores {
		fmt.Printf("key=%s, value=%d\n", k, v)
	}

	fmt.Println("\n--- range 遍历字符串（按 rune）---")
	for i, r := range "Hello, 世界" {
		fmt.Printf("index=%d, rune=%c\n", i, r)
	}

	fmt.Println("=========start to show cap demo ==========")
	test1()
}

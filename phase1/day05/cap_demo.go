package main

import "fmt"

func test1() {
	// ============================================================
	//  切片 (slice) 的底层结构：
	//  ┌─────────────┬─────┬─────┐
	//  │   指针 ptr   │ len │ cap │
	//  └─────────────┴─────┴─────┘
	//
	//  - ptr: 指向底层数组的起始位置
	//  - len: 当前能访问的元素个数（s[0] ~ s[len-1]）
	//  - cap: 从 ptr 指向的位置到底层数组末尾的总容量
	// ============================================================

	fmt.Println("=== 1. 从数组切片：cap 怎么算？ ===")
	arr := [6]int{0, 10, 20, 30, 40, 50}
	// 数组内存布局（连续）：
	// index:  0    1    2    3    4    5
	// value: [0] [10] [20] [30] [40] [50]

	s1 := arr[1:4] // 取 index 1,2,3 → 值 [10,20,30]
	// s1.ptr 指向 arr[1]
	// s1.len = 4-1 = 3
	// s1.cap = 从 arr[1] 到数组末尾 = 6-1 = 5
	fmt.Printf("arr = %v\n", arr)
	fmt.Printf("s1 = arr[1:4] = %v, len=%d, cap=%d\n", s1, len(s1), cap(s1))
	fmt.Println("解释: cap=5 是因为从 arr[1] 到数组末尾还有 5 个位置(1~5)")

	fmt.Println("\n=== 2. 对切片再切片：cap 继承 ===")
	s2 := s1[1:3] // 从 s1 的 index 1~2 取 → 值 [20,30]
	// s2.ptr 指向 s1[1] 也就是 arr[2]
	// s2.len = 3-1 = 2
	// s2.cap = 从 arr[2] 到数组末尾 = 6-2 = 4
	fmt.Printf("s2 = s1[1:3] = %v, len=%d, cap=%d\n", s2, len(s2), cap(s2))
	fmt.Println("解释: 虽然 s1.len=3，但 s2.cap=4 是因为底层数组还剩4个位置")

	fmt.Println("\n=== 3. 超出 len 但在 cap 范围内：可以访问！ ===")
	// s1 = [10,20,30], len=3, cap=5
	// 这意味着 s1 底层还有 arr[4]=40, arr[5]=50 可以"看到"
	fmt.Printf("s1 的值: %v (len=%d)\n", s1[:cap(s1)], cap(s1))
	fmt.Printf("用切片语法 s1[:cap(s1)] 可以看到: %v\n", s1[:cap(s1)])
	fmt.Println("⚠️ 危险！这样可以访问到 arr[4], arr[5]，因为 cap 允许！")

	fmt.Println("\n=== 4. make([]T, len, cap)：预分配容量 ===")
	s3 := make([]int, 2, 6)
	s3[0], s3[1] = 100, 200
	fmt.Printf("make 后: s3=%v, len=%d, cap=%d\n", s3, len(s3), cap(s3))

	// 前 4 次 append 不会扩容，因为 cap=6
	for i := 1; i <= 4; i++ {
		oldCap := cap(s3)
		s3 = append(s3, i*10)
		if cap(s3) != oldCap {
			fmt.Printf("第%d次 append 扩容了! cap %d -> %d\n", i, oldCap, cap(s3))
		} else {
			fmt.Printf("第%d次 append: s3=%v, len=%d, cap=%d (未扩容)\n", i, s3, len(s3), cap(s3))
		}
	}

	fmt.Println("\n=== 5. append 扩容规则：超过 cap 才会扩容 ===")
	s4 := make([]int, 2, 2)
	s4[0], s4[1] = 1, 2
	fmt.Printf("初始: s4=%v, len=%d, cap=%d, ptr=%p\n", s4, len(s4), cap(s4), s4)

	// 第1次 append: len=3 > cap=2，触发扩容
	oldPtr := &s4[0]
	s4 = append(s4, 3)
	fmt.Printf("append 3 后: s4=%v, len=%d, cap=%d, ptr=%p\n", s4, len(s4), cap(s4), s4)
	fmt.Printf("底层数组变了? %v (指针不同说明分配了新数组)\n", oldPtr != &s4[0])

	fmt.Println("\n=== 6. append 扩容时 cap 的增长策略 ===")
	sizes := []int{}
	s := make([]int, 0)
	for i := 0; i < 20; i++ {
		if cap(s) != 0 {
			last := sizes[len(sizes)-1]
			if cap(s) != last {
				sizes = append(sizes, cap(s))
				fmt.Printf("append 第%2d个元素后: len=%2d, cap=%2d (扩容, 系数≈%.2f)\n",
					i, len(s), cap(s), float64(cap(s))/float64(last))
			}
		} else {
			sizes = append(sizes, cap(s))
		}
		s = append(s, i)
	}
	fmt.Println("\n扩容规律:")
	fmt.Println("  - 小于 1024: 容量翻倍 (2x)")
	fmt.Println("  - 大于 1024: 大约 1.25x 增长")

	fmt.Println("\n=== 7. 为什么关注 cap？性能陷阱！ ===")
	fmt.Println("❌ 低效写法（频繁扩容，多次内存分配+数据拷贝）:")
	// bad := []int{}
	// for i := 0; i < 100000; i++ {
	//     bad = append(bad, i)  // 扩容 10+ 次
	// }

	fmt.Println("✅ 高效写法（一次性预分配）:")
	// good := make([]int, 0, 100000)
	// for i := 0; i < 100000; i++ {
	//     good = append(good, i)  // 0 次扩容
	// }

	fmt.Println("\n=== 8. 切片传参：cap 也一起传 ===")
	modifySlice := func(s []int) {
		// 这里接收的是切片的拷贝（头信息拷贝，ptr/len/cap 都拷贝）
		// 但 ptr 指向同一个底层数组，所以修改会影响原切片
		s[0] = 999
		// append 如果未超过 cap，也会影响原底层数组
		_ = append(s, 888)
		fmt.Printf("函数内 s=%v, len=%d, cap=%d\n", s, len(s), cap(s))
	}

	orig := []int{1, 2, 3, 4, 5}
	fmt.Printf("调用前 orig=%v\n", orig)
	modifySlice(orig[:3]) // 只传前3个，但 cap 还是 5
	fmt.Printf("调用后 orig=%v  ← 底层数组被修改了！\n", orig)
}

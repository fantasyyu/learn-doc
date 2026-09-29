package main

import (
	"fmt"
	"strings"
	"time"
)

// === 1. 数据结构定义 ===

// Book 图书
type Book struct {
	ID       int
	Title    string
	Author   string
	ISBN     string
	Borrowed bool
	BorrowAt *time.Time
}

// Member 借阅者
type Member struct {
	ID   int
	Name string
}

// === 2. 错误定义 ===

type LibError struct {
	Op  string // 操作
	Msg string
}

func (e *LibError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Op, e.Msg)
}

// === 3. 图书馆（核心管理结构）===
type Library struct {
	books  []Book
	members []Member
	nextBookID   int
	nextMemberID int
}

// NewLibrary 创建图书馆
func NewLibrary() *Library {
	return &Library{
		books:        make([]Book, 0),
		members:      make([]Member, 0),
		nextBookID:   1,
		nextMemberID: 1,
	}
}

// AddBook 添加图书（指针接收者：修改内部状态）
func (l *Library) AddBook(title, author, isbn string) *Book {
	book := Book{
		ID:     l.nextBookID,
		Title:  title,
		Author: author,
		ISBN:   isbn,
	}
	l.nextBookID++
	l.books = append(l.books, book)
	return &l.books[len(l.books)-1]
}

// RegisterMember 注册会员
func (l *Library) RegisterMember(name string) *Member {
	m := Member{ID: l.nextMemberID, Name: name}
	l.nextMemberID++
	l.members = append(l.members, m)
	return &l.members[len(l.members)-1]
}

// FindBookByID 按 ID 查找
func (l *Library) FindBookByID(id int) (*Book, error) {
	for i := range l.books {
		if l.books[i].ID == id {
			return &l.books[i], nil
		}
	}
	return nil, &LibError{Op: "FindBook", Msg: fmt.Sprintf("图书 #%d 不存在", id)}
}

// SearchByTitle 按标题搜索（支持模糊匹配）
func (l *Library) SearchByTitle(keyword string) []Book {
	results := make([]Book, 0)
	for _, b := range l.books {
		if strings.Contains(strings.ToLower(b.Title), strings.ToLower(keyword)) {
			results = append(results, b)
		}
	}
	return results
}

// Borrow 借阅
func (l *Library) Borrow(bookID int, memberID int) error {
	book, err := l.FindBookByID(bookID)
	if err != nil {
		return err
	}
	if book.Borrowed {
		return &LibError{Op: "Borrow", Msg: fmt.Sprintf("《%s》已被借出", book.Title)}
	}
	now := time.Now()
	book.Borrowed = true
	book.BorrowAt = &now
	fmt.Printf("✓ %s 借阅了《%s》\n", l.memberName(memberID), book.Title)
	return nil
}

// Return 归还
func (l *Library) Return(bookID int) error {
	book, err := l.FindBookByID(bookID)
	if err != nil {
		return err
	}
	if !book.Borrowed {
		return &LibError{Op: "Return", Msg: fmt.Sprintf("《%s》未被借出", book.Title)}
	}
	book.Borrowed = false
	book.BorrowAt = nil
	fmt.Printf("✓ 《%s》已归还\n", book.Title)
	return nil
}

// ListBooks 列出所有图书
func (l *Library) ListBooks() {
	fmt.Println("\n📚 图书列表：")
	fmt.Printf("%-4s %-20s %-15s %-13s %s\n", "ID", "书名", "作者", "ISBN", "状态")
	fmt.Println(strings.Repeat("-", 70))
	for _, b := range l.books {
		status := "可借阅"
		if b.Borrowed {
			status = fmt.Sprintf("已借出 (%s)", b.BorrowAt.Format("01-02"))
		}
		fmt.Printf("%-4d %-20s %-15s %-13s %s\n", b.ID, b.Title, b.Author, b.ISBN, status)
	}
}

// 辅助方法
func (l *Library) memberName(id int) string {
	for _, m := range l.members {
		if m.ID == id {
			return m.Name
		}
	}
	return "未知用户"
}

// === 4. 主程序 ===
func main() {
	lib := NewLibrary()

	// 注册会员
	alice := lib.RegisterMember("Alice")
	bob := lib.RegisterMember("Bob")
	fmt.Printf("注册会员: %s(ID=%d), %s(ID=%d)\n", alice.Name, alice.ID, bob.Name, bob.ID)

	// 添加图书
	lib.AddBook("Go程序设计语言", "Alan Donovan", "978-0134190440")
	lib.AddBook("深入理解计算机系统", "Randal Bryant", "978-0134092669")
	lib.AddBook("设计模式", "Erich Gamma", "978-0201633610")
	lib.AddBook("算法导论", "Thomas Cormen", "978-0262033848")
	fmt.Println("已添加 4 本图书")

	// 列出
	lib.ListBooks()

	// 借阅
	fmt.Println("\n--- 借阅操作 ---")
	_ = lib.Borrow(1, alice.ID) // Alice 借《Go程序设计语言》
	_ = lib.Borrow(2, bob.ID)   // Bob 借《深入理解计算机系统》
	_ = lib.Borrow(1, bob.ID)   // 重复借：应失败

	// 搜索
	fmt.Println("\n--- 搜索 '设计' ---")
	results := lib.SearchByTitle("设计")
	for _, b := range results {
		fmt.Printf("找到: 《%s》 by %s\n", b.Title, b.Author)
	}

	// 归还
	fmt.Println("\n--- 归还操作 ---")
	_ = lib.Return(1)
	_ = lib.Return(99) // 不存在的书

	// 最终列表
	lib.ListBooks()

	fmt.Println("\n✅ 图书管理系统运行完毕！")
}

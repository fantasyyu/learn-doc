// GORM 完整 CRUD + 关联关系演示
// 使用 SQLite（文件数据库），无需安装任何服务！
// 换 MySQL 只需改 Open 那一行：
//   dsn := "root:123456@tcp(127.0.0.1:3306)/godemo?charset=utf8mb4&parseTime=True"
//   db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"

	"github.com/glebarez/sqlite" // 纯 Go 实现，无需 CGO（生产环境可换 gorm.io/driver/mysql）
)

// ============================================================
//  模型定义 — GORM 通过约定自动识别表结构
// ============================================================

// User 主表
// 约定: 结构体名=表名(User->users), ID字段=主键, CreatedAt/UpdatedAt自动维护
type User struct {
	ID        uint           `gorm:"primaryKey"`
	Name      string         `gorm:"size:100;not null;index"`
	Email     string         `gorm:"size:100;uniqueIndex"` // uniqueIndex: 唯一索引
	Age       int
	CreditCard CreditCard    // 一对一关联（拥有）
	Orders    []Order        // 一对多关联
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreditCard 一对一（User 拥有 CreditCard）
type CreditCard struct {
	ID     uint
	Number string
	UserID uint // 外键：关联到 User.ID
}

// Order 一对多（User 有多个 Order）
type Order struct {
	ID     uint
	Amount float64
	UserID uint // 外键
}

// 表名自定义（可选，默认是结构体名的蛇形复数）
func (User) TableName() string {
	return "my_users"
}

func main() {
	// ============================================================
	//  1. 连接 + 自动迁移（AutoMigrate 会按模型建表/补列）
	// ============================================================
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("连接失败:", err)
	}
	fmt.Println("✓ 数据库已连接 (SQLite 文件: test.db)")

	// AutoMigrate: 自动创建/更新表结构（生产环境建议用专门的 migration 工具）
	db.AutoMigrate(&User{}, &CreditCard{}, &Order{})
	fmt.Println("✓ 表结构自动迁移完成")

	// ============================================================
	//  2. 创建 (Create)
	// ============================================================
	user := User{Name: "Alice", Email: "alice@example.com", Age: 25,
		CreditCard: CreditCard{Number: "6222-0000-1111"},
		Orders:     []Order{{Amount: 99.5}, {Amount: 199.0}},
	}
	result := db.Create(&user) // 关联数据会一并插入！
	fmt.Printf("✓ 创建用户: ID=%d, 行数=%d\n", user.ID, result.RowsAffected)

	// 批量创建
	users := []User{
		{Name: "Bob", Email: "bob@example.com", Age: 30},
		{Name: "Charlie", Email: "charlie@example.com", Age: 35},
	}
	db.CreateInBatches(users, 100)
	fmt.Println("✓ 批量创建 2 个用户")

	// ============================================================
	//  3. 查询 (First / Find / Where)
	// ============================================================
	var u User
	db.First(&u, 1) // 按主键查
	fmt.Printf("✓ First(主键=1): %s\n", u.Name)

	db.First(&u, "email = ?", "bob@example.com") // 条件查询
	fmt.Printf("✓ Where(email): %s\n", u.Name)

	var all []User
	db.Find(&all) // 查全部
	fmt.Printf("✓ Find 全部: %d 个用户\n", len(all))

	// 条件链式 + 排序 + 分页
	var adults []User
	db.Where("age >= ?", 18).Order("age desc").Limit(10).Offset(0).Find(&adults)
	for _, a := range adults {
		fmt.Printf("    成人用户: %s (%d岁)\n", a.Name, a.Age)
	}

	// ============================================================
	//  4. 预加载关联 (Preload)
	// ============================================================
	var userWithOrders User
	db.Preload("Orders").Preload("CreditCard").First(&userWithOrders, 1)
	fmt.Printf("✓ 预加载: %s, 订单数=%d, 卡号=%s\n",
		userWithOrders.Name, len(userWithOrders.Orders), userWithOrders.CreditCard.Number)

	// ============================================================
	//  5. 更新 (Update / Updates)
	// ============================================================
	db.Model(&User{}).Where("id = ?", 1).Update("age", 26)          // 更新单列
	db.Model(&User{}).Where("id = ?", 1).Updates(map[string]interface{}{ // 多列
		"name": "Alice Cooper",
		"age":  27,
	})
	fmt.Println("✓ 更新完成")

	// ============================================================
	//  6. 删除 (Delete)
	// ============================================================
	db.Delete(&User{}, 3) // 按主键删除 Charlie
	fmt.Println("✓ 删除 ID=3")

	// ============================================================
	//  7. 事务
	// ============================================================
	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&User{Name: "Dave", Email: "dave@example.com", Age: 40}).Error; err != nil {
			return err // 返回 error 自动回滚
		}
		if err := tx.Create(&Order{Amount: 500, UserID: 4}).Error; err != nil {
			return err
		}
		return nil // 返回 nil 提交
	})
	if err != nil {
		log.Fatal("事务失败:", err)
	}
	fmt.Println("✓ 事务提交成功")

	// ============================================================
	//  8. 统计
	// ============================================================
	var count int64
	db.Model(&User{}).Where("age >= ?", 18).Count(&count)
	fmt.Printf("✓ 成人总数: %d\n", count)

	fmt.Println("\n=== GORM vs database/sql ===")
	fmt.Println("database/sql — 手写 SQL，灵活但要自己处理映射")
	fmt.Println("GORM        — 链式 API，CRUD+关联一把梭，复杂 SQL 用 Raw()")
}

// MySQL 标准库操作完整示例
// 运行前提：本地有 MySQL 服务，且已创建数据库：
//   CREATE DATABASE godemo CHARACTER SET utf8mb4;
//
// 启动方式：
//   go run main.go
// 或先用 Docker 启动 MySQL：
//   docker run -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:8
package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // 匿名导入：只执行其 init() 注册驱动
)

// User 对应 users 表
type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

func main() {
	// ============================================================
	//  1. 连接数据库
	//  DSN 格式: 用户名:密码@tcp(地址:端口)/数据库名?参数
	// ============================================================
	dsn := "root:123456@tcp(127.0.0.1:3306)/godemo?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := sql.Open("mysql", dsn) // Open 只是验证参数，不真正连接
	if err != nil {
		log.Fatal("DSN 解析失败:", err)
	}
	defer db.Close()

	// 配置连接池（生产必配！）
	db.SetMaxOpenConns(25)                 // 最大打开连接数
	db.SetMaxIdleConns(10)                 // 最大空闲连接数
	db.SetConnMaxLifetime(5 * time.Minute) // 连接最大存活时间

	// Ping 才真正建立连接
	if err := db.Ping(); err != nil {
		log.Fatal("连接 MySQL 失败（请先启动 MySQL 并创建 godemo 库）:", err)
	}
	fmt.Println("✓ MySQL 连接成功")

	// ============================================================
	//  2. 建表
	// ============================================================
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal("建表失败:", err)
	}
	fmt.Println("✓ users 表就绪")

	// ============================================================
	//  3. 插入 — Exec + LastInsertId
	// ============================================================
	// 占位符用 ?（MySQL），永远不要用字符串拼接（防 SQL 注入！）
	result, err := db.Exec(
		"INSERT INTO users (name, email) VALUES (?, ?)",
		"Alice", "alice@example.com",
	)
	if err != nil {
		log.Fatal("插入失败:", err)
	}
	lastID, _ := result.LastInsertId()
	fmt.Printf("✓ 插入成功, ID=%d\n", lastID)

	// 批量插入（预处理语句，性能好且安全）
	stmt, _ := db.Prepare("INSERT INTO users (name, email) VALUES (?, ?)")
	defer stmt.Close()
	for _, u := range []User{{Name: "Bob", Email: "bob@example.com"}, {Name: "Charlie", Email: "charlie@example.com"}} {
		stmt.Exec(u.Name, u.Email)
	}
	fmt.Println("✓ 批量插入 2 条")

	// ============================================================
	//  4. 单行查询 — QueryRow
	// ============================================================
	var user User
	err = db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", lastID).
		Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt) // Scan 顺序必须与 SELECT 一致
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	fmt.Printf("✓ 查询单行: %+v\n", user)

	// ============================================================
	//  5. 多行查询 — Query + rows.Next
	// ============================================================
	rows, err := db.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		log.Fatal("查询失败:", err)
	}
	defer rows.Close() // 必须关闭！否则连接泄漏

	fmt.Println("✓ 查询多行:")
	for rows.Next() {
		var u User
		rows.Scan(&u.ID, &u.Name, &u.Email)
		fmt.Printf("  #%d %s <%s>\n", u.ID, u.Name, u.Email)
	}
	if err := rows.Err(); err != nil { // 遍历结束后检查错误
		log.Fatal("遍历出错:", err)
	}

	// ============================================================
	//  6. 事务 — 保证多步操作的原子性
	// ============================================================
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("开启事务失败:", err)
	}

	// 两步操作，要么都成功，要么都回滚
	if _, err := tx.Exec("UPDATE users SET name = ? WHERE id = ?", "Alice Cooper", lastID); err != nil {
		tx.Rollback()
		log.Fatal("事务回滚:", err)
	}
	if _, err := tx.Exec("UPDATE users SET email = ? WHERE id = ?", "alice@cooper.com", lastID); err != nil {
		tx.Rollback()
		log.Fatal("事务回滚:", err)
	}
	if err := tx.Commit(); err != nil { // 提交事务
		tx.Rollback()
		log.Fatal("提交失败:", err)
	}
	fmt.Println("✓ 事务提交成功（两条 UPDATE 已生效）")

	// ============================================================
	//  7. 删除
	// ============================================================
	res, _ := db.Exec("DELETE FROM users WHERE name = ?", "Charlie")
	affected, _ := res.RowsAffected()
	fmt.Printf("✓ 删除 %d 行\n", affected)

	fmt.Println("\n=== database/sql 要点 ===")
	fmt.Println("1. sql.Open 不真正连接，Ping 才连接")
	fmt.Println("2. 永远用 ? 占位符，不要拼字符串（防注入）")
	fmt.Println("3. rows 必须 Close，连接池才不会耗尽")
	fmt.Println("4. 多步写操作用事务 Begin/Commit/Rollback")
	fmt.Println("5. 连接池三参数：MaxOpen/MaxIdle/ConnMaxLifetime")
}

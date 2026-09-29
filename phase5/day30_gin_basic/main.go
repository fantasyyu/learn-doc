package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// User 请求体绑定的结构体
type User struct {
	Name string `json:"name" binding:"required"`  // binding: 校验标签
	Age  int    `json:"age" binding:"gte=0,lte=150"` // gte/lte: 范围校验
}

// 全局模拟数据
var users = []map[string]interface{}{
	{"id": 1, "name": "Alice"},
	{"id": 2, "name": "Bob"},
}

// 中间件：统计请求耗时
func latencyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next() // 继续执行后续处理器
		latency := time.Since(start)
		c.Header("X-Latency", latency.String())
	}
}

func main() {
	// 默认模式（带 Logger 和 Recovery 中间件）
	r := gin.Default()

	// 注册全局中间件
	r.Use(latencyMiddleware())

	// ============================================================
	//  1. 基本路由
	// ============================================================
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// ============================================================
	//  2. 路径参数 /users/123
	// ============================================================
	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id") // 取路径参数
		c.JSON(http.StatusOK, gin.H{
			"user_id": id,
		})
	})

	// ============================================================
	//  3. Query 参数 /search?q=golang&page=2
	// ============================================================
	r.GET("/search", func(c *gin.Context) {
		q := c.DefaultQuery("q", "")          // 带默认值
		page := c.DefaultQuery("page", "1")   // 带默认值
		c.JSON(http.StatusOK, gin.H{
			"query": q,
			"page":  page,
		})
	})

	// ============================================================
	//  4. 请求体 JSON 绑定 + 校验
	// ============================================================
	r.POST("/users", func(c *gin.Context) {
		var user User
		// ShouldBindJSON: 解析 JSON 并校验 binding 标签
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"created": user,
		})
	})

	// ============================================================
	//  5. 返回全部用户
	// ============================================================
	r.GET("/users", func(c *gin.Context) {
		c.JSON(http.StatusOK, users)
	})

	// ============================================================
	//  6. 路由组 — API 版本控制
	// ============================================================
	v1 := r.Group("/api/v1")
	{
		v1.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"version": "v1"})
		})
	}
	v2 := r.Group("/api/v2")
	{
		v2.GET("/info", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"version": "v2"})
		})
	}

	// ============================================================
	//  7. 路由参数通配符
	// ============================================================
	r.GET("/files/*filepath", func(c *gin.Context) {
		path := c.Param("filepath") // 匹配 /files/a/b/c.txt
		c.JSON(http.StatusOK, gin.H{"path": path})
	})

	// 启动服务器
	r.Run(":8080") // 默认监听 0.0.0.0:8080
}

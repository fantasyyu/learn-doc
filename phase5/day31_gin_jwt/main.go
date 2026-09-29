package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// 密钥（生产环境必须放在环境变量/配置中心！）
var jwtSecret = []byte("my-secret-key-change-me-in-production")

// Claims — 自定义 JWT 负载
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// ============================================================
//  生成 JWT Token
// ============================================================
func generateToken(userID int, role string) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)), // 2小时过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "my-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ============================================================
//  JWT 认证中间件
// ============================================================
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing or invalid Authorization header",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析并验证 token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			// 验证签名算法，防止算法混淆攻击
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token: " + err.Error(),
			})
			return
		}

		// 把用户信息存入上下文，供后续处理器使用
		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Set("role", claims.Role)
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
			})
		}
	}
}

// ============================================================
//  角色鉴权中间件（RBAC）
// ============================================================
func requireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists || userRole.(string) != role {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}
		c.Next()
	}
}

// ============================================================
//  统一错误处理中间件
// ============================================================
func errorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// 捕获处理过程中产生的错误
		if len(c.Errors) > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errors": c.Errors.Errors(),
			})
		}
	}
}

func main() {
	r := gin.Default()
	r.Use(errorHandler())

	// ============ 公开路由 ============
	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 模拟验证（实际项目查数据库）
		var userID int
		var role string
		switch req.Username {
		case "admin":
			userID, role = 1, "admin"
		case "alice":
			userID, role = 2, "user"
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}

		token, err := generateToken(userID, role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"token": token,
			"role":  role,
			"note":  "2小时后过期",
		})
	})

	// ============ 需要认证的路由组 ============
	authorized := r.Group("/api", authMiddleware())
	{
		// 任何登录用户都能访问
		authorized.GET("/profile", func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")
			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"role":    role,
			})
		})

		// 只有 admin 能访问
		authorized.GET("/admin/dashboard", requireRole("admin"), func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "欢迎来到管理后台",
				"stats":   gin.H{"users": 100, "orders": 520},
			})
		})
	}

	r.Run(":8081")
}

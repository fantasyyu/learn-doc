// viper 配置管理演示
// 支持: YAML/JSON/TOML/环境变量/命令行/远程配置(K8s ConfigMap等)
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// Config — 把配置映射到结构体（推荐：类型安全 + IDE 补全）
// mapstructure 标签解决蛇形命名配置 -> 驼峰字段的映射问题
type Config struct {
	Server struct {
		Port int
		Mode string
	}
	Database struct {
		Driver       string
		DSN          string `mapstructure:"dsn"`
		MaxOpenConns int    `mapstructure:"max_open_conns"`
		MaxIdleConns int    `mapstructure:"max_idle_conns"`
	}
	Redis struct {
		Addr     string
		Password string
		DB       int
	}
	JWT struct {
		Secret      string
		ExpireHours int `mapstructure:"expire_hours"`
	}
	Log struct {
		Level  string
		Format string
	}
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	Features       map[string]interface{}
}

func main() {
	// ============================================================
	//  1. 基础配置加载
	// ============================================================
	viper.SetConfigName("config")   // 文件名（不含扩展名）
	viper.SetConfigType("yaml")     // 配置格式
	viper.AddConfigPath(".")        // 查找路径
	viper.AddConfigPath("./config") // 可以放 config 子目录

	// 默认值：配置文件缺失时兜底
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("log.level", "info")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("读取配置失败: %w", err))
	}
	fmt.Printf("✓ 已加载配置文件: %s\n", viper.ConfigFileUsed())

	// ============================================================
	//  2. 读取配置值的三种方式
	// ============================================================
	fmt.Println("\n=== 读取配置 ===")

	// 方式A: Get（返回 interface{}，需要类型断言）
	port := viper.GetInt("server.port")
	fmt.Printf("server.port: %d\n", port)

	// 方式B: 嵌套访问
	fmt.Printf("database.driver: %s\n", viper.GetString("database.driver"))

	// 方式C: Unmarshal 到结构体（推荐！）
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	fmt.Printf("JWT 密钥: %s (过期: %d小时)\n", cfg.JWT.Secret, cfg.JWT.ExpireHours)
	fmt.Printf("允许的来源: %v\n", cfg.AllowedOrigins)
	fmt.Printf("特性开关: 缓存=%v, 限流=%v\n",
		cfg.Features["enable_cache"], cfg.Features["rate_limit"])

	// ============================================================
	//  3. 环境变量覆盖（12-Factor 原则）
	//  生产环境用环境变量覆盖配置文件，无需改文件
	// ============================================================
	fmt.Println("\n=== 环境变量覆盖 ===")

	// 设置环境变量前缀，并自动把 SERVER_PORT 映射到 server.port
	viper.SetEnvPrefix("APP")                 // APP_SERVER_PORT -> ...
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()                      // 自动读取环境变量

	// 模拟设置环境变量后读取
	fmt.Println("在终端执行: export APP_SERVER_PORT=9090")
	fmt.Println("然后 viper.GetInt(server.port) 会返回 9090")
	fmt.Printf("当前 server.port = %d\n", viper.GetInt("server.port"))

	// ============================================================
	//  4. 多配置文件（开发/测试/生产分离）
	// ============================================================
	fmt.Println("\n=== 多环境配置 ===")
	fmt.Println("viper 支持按环境加载不同配置:")
	fmt.Println("  config.yaml         — 基础配置")
	fmt.Println("  config.dev.yaml     — 开发环境覆盖")
	fmt.Println("  config.prod.yaml    — 生产环境覆盖")
	fmt.Println("")
	fmt.Println("常见做法:")
	fmt.Println("  env := os.Getenv(APP_ENV)  // dev / prod")
	fmt.Println("  viper.SetConfigName(config. + env)")
	fmt.Println("  viper.MergeInConfig()      // 合并覆盖")

	// ============================================================
	//  5. 配置文件热更新（监听变化）
	// ============================================================
	fmt.Println("\n=== 配置热更新 ===")
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("→ 配置文件被修改，自动重载！")
	})
	// viper.WatchConfig() // 取消注释启用（演示中不阻塞）
	fmt.Println("viper.WatchConfig() 可监听文件变化自动重载")

	// ============================================================
	//  6. 实战：把配置用于连接信息打印
	// ============================================================
	fmt.Println("\n=== 启动信息汇总 ===")
	fmt.Printf("服务将在端口 %d 启动（模式: %s）\n", cfg.Server.Port, cfg.Server.Mode)
	fmt.Printf("数据库: %s\n", maskDSN(cfg.Database.DSN))
	fmt.Printf("Redis: %s (db=%d)\n", cfg.Redis.Addr, cfg.Redis.DB)
	fmt.Printf("日志级别: %s, 格式: %s\n", cfg.Log.Level, cfg.Log.Format)
	fmt.Printf("Token 过期时间: %d小时\n", cfg.JWT.ExpireHours)

	time.Sleep(100 * time.Millisecond)
	fmt.Println("\n✓ viper 是 Go 生态配置管理的事实标准")
}

// maskDSN — 隐藏密码，防止日志泄露
func maskDSN(dsn string) string {
	if at := strings.Index(dsn, "@"); at > 0 {
		if colon := strings.Index(dsn[:at], ":"); colon > 0 {
			return dsn[:colon+1] + "***" + dsn[at:]
		}
	}
	return dsn
}

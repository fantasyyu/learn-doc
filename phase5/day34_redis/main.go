// Redis 完整操作演示
// 使用 miniredis（内存版 Redis）演示，无需安装 Redis 服务！
// 连接真实 Redis 只需把下面的 miniredis 换成：
//   rdb := redis.NewClient(&redis.Options{
//       Addr:     "localhost:6379",
//       Password: "", // 有密码就填
//       DB:       0,
//   })
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	// ============================================================
	//  启动内存版 Redis（仅演示用；生产环境连真实 Redis）
	// ============================================================
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}
	defer mr.Close()

	// 连接（Addr 指向 miniredis；生产环境改成 localhost:6379）
	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	defer rdb.Close()

	// Ping 测试
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic(err)
	}
	fmt.Printf("✓ Redis 已连接 (%s)\n", mr.Addr())

	// ============================================================
	//  1. String — 最常用：缓存、计数器、验证码
	// ============================================================
	fmt.Println("\n=== 1. String ===")

	rdb.Set(ctx, "name", "Alice", 1*time.Hour) // 带过期时间
	rdb.Set(ctx, "temp", "xyz", 0)             // 0 = 永不过期

	name, _ := rdb.Get(ctx, "name").Result()
	fmt.Printf("GET name = %s\n", name)

	// 自增/自减（计数器场景：点赞数、库存）
	rdb.Set(ctx, "likes", 0, 0)
	rdb.Incr(ctx, "likes")     // +1
	rdb.IncrBy(ctx, "likes", 5) // +5
	likes, _ := rdb.Get(ctx, "likes").Result()
	fmt.Printf("点赞数: %s\n", likes)

	// 批量操作
	rdb.MSet(ctx, "k1", "v1", "k2", "v2")
	vals, _ := rdb.MGet(ctx, "k1", "k2", "不存在").Result()
	fmt.Printf("MGET: %v ← 不存在的key返回nil\n", vals)

	// 过期时间
	ttl, _ := rdb.TTL(ctx, "name").Result()
	fmt.Printf("name 的 TTL: %v\n", ttl)
	rdb.Expire(ctx, "temp", 30*time.Minute) // 动态设置过期

	// ============================================================
	//  2. Hash — 存储对象（比 String 省内存，可单独改字段）
	// ============================================================
	fmt.Println("\n=== 2. Hash ===")

	rdb.HSet(ctx, "user:1", "name", "Bob", "age", "30", "city", "Beijing")
	rdb.HSet(ctx, "user:1", "age", "31") // 单独更新一个字段

	age, _ := rdb.HGet(ctx, "user:1", "age").Result()
	fmt.Printf("user:1.age = %s\n", age)

	userFields, _ := rdb.HGetAll(ctx, "user:1").Result()
	fmt.Printf("user:1 全部字段: %v\n", userFields)

	// Hash 适合做对象缓存，比如缓存数据库用户：
	fmt.Println("  典型用法: 查库前先 HGETALL user:{id}， miss 再查库并缓存")

	// ============================================================
	//  3. List — 队列/最新消息
	// ============================================================
	fmt.Println("\n=== 3. List ===")

	rdb.LPush(ctx, "news", "新闻3", "新闻2", "新闻1") // 左进（最新在前）
	rdb.RPush(ctx, "queue", "任务A", "任务B")          // 右进（排队）

	// 最新 2 条新闻
	news, _ := rdb.LRange(ctx, "news", 0, 1).Result()
	fmt.Printf("最新新闻: %v\n", news)

	// 队列：左边出 = FIFO
	task, _ := rdb.LPop(ctx, "queue").Result()
	fmt.Printf("取出任务: %s (剩余队列长度=%d)\n", task, rdb.LLen(ctx, "queue").Val())

	// ============================================================
	//  4. Set — 去重、标签、共同好友
	// ============================================================
	fmt.Println("\n=== 4. Set ===")

	rdb.SAdd(ctx, "tags:go", "后端", "云原生", "高性能")
	rdb.SAdd(ctx, "tags:python", "后端", "AI", "数据分析")

	common, _ := rdb.SInter(ctx, "tags:go", "tags:python").Result() // 交集
	fmt.Printf("Go和Python的共同标签: %v\n", common)

	isMember, _ := rdb.SIsMember(ctx, "tags:go", "云原生").Result()
	fmt.Printf("go 有云原生标签? %v\n", isMember)

	// 场景：防止重复提交
	rdb.SAdd(ctx, "submitted:1001", "order-123")
	ok, _ := rdb.SAdd(ctx, "submitted:1001", "order-123").Result()
	fmt.Printf("重复提交检测: order-123 重复提交? %v (返回0=已存在)\n", ok == 0)

	// ============================================================
	//  5. Sorted Set — 排行榜（最经典场景）
	// ============================================================
	fmt.Println("\n=== 5. Sorted Set (排行榜) ===")

	// 添加玩家分数
	rdb.ZAdd(ctx, "game:rank", redis.Z{Score: 100, Member: "玩家A"},
		redis.Z{Score: 200, Member: "玩家B"},
		redis.Z{Score: 150, Member: "玩家C"})

	// 更新分数
	rdb.ZIncrBy(ctx, "game:rank", 80, "玩家A")

	// 取前 3 名（分数从高到低）
	rank, _ := rdb.ZRevRangeWithScores(ctx, "game:rank", 0, 2).Result()
	fmt.Println("🏆 排行榜 Top3:")
	for i, z := range rank {
		fmt.Printf("  第%d名: %s (%.0f分)\n", i+1, z.Member, z.Score)
	}

	// 查某个玩家的排名
	pos, _ := rdb.ZRevRank(ctx, "game:rank", "玩家A").Result()
	fmt.Printf("玩家A 的排名: 第%d名\n", pos+1)

	// ============================================================
	//  6. 分布式锁（SET NX EX）
	// ============================================================
	fmt.Println("\n=== 6. 分布式锁 ===")

	// SET key value NX EX seconds — 不存在才设置，带过期（防死锁）
	lockKey := "lock:order"
	lockVal := "server-1"
	acquired, _ := rdb.SetNX(ctx, lockKey, lockVal, 10*time.Second).Result()
	fmt.Printf("获取锁: %v\n", acquired)

	// 再次获取（已被占用）
	acquired2, _ := rdb.SetNX(ctx, lockKey, "server-2", 10*time.Second).Result()
	fmt.Printf("其他服务获取同一把锁: %v ← 被占用，获取失败\n", acquired2)

	// 释放锁（先确认是自己的锁再删，防止误删别人的锁）
	if v, _ := rdb.Get(ctx, lockKey).Result(); v == lockVal {
		rdb.Del(ctx, lockKey)
		fmt.Println("✓ 释放锁")
	}

	// ============================================================
	//  7. Pipeline — 批量操作减少网络往返
	// ============================================================
	fmt.Println("\n=== 7. Pipeline ===")

	start := time.Now()
	pipe := rdb.Pipeline()
	for i := 0; i < 100; i++ {
		pipe.Set(ctx, fmt.Sprintf("key:%d", i), i, time.Hour)
	}
	cmds, _ := pipe.Exec(ctx) // 一次网络往返执行100条命令
	fmt.Printf("Pipeline 执行 %d 条命令，耗时: %v\n", len(cmds), time.Since(start))

	fmt.Println("\n=== 总结 ===")
	fmt.Println("String       — 缓存、计数器、分布式锁")
	fmt.Println("Hash         — 对象存储（可单字段更新）")
	fmt.Println("List         — 队列、最新消息")
	fmt.Println("Set          — 去重、标签交集")
	fmt.Println("Sorted Set   — 排行榜")
	fmt.Println("Pipeline     — 批量操作提性能")
}

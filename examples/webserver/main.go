package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/scache/cache"
)

// User 用户数据结构
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// 模拟数据库
var mockDB = map[int]User{
	1: {ID: 1, Name: "张三", Email: "zhangsan@example.com", Age: 25},
	2: {ID: 2, Name: "李四", Email: "lisi@example.com", Age: 30},
	3: {ID: 3, Name: "王五", Email: "wangwu@example.com", Age: 28},
	4: {ID: 4, Name: "赵六", Email: "zhaoliu@example.com", Age: 35},
	5: {ID: 5, Name: "钱七", Email: "qianqi@example.com", Age: 22},
}

// 模拟从数据库加载用户
func loadUserFromDB(id int) (User, error) {
	// 模拟数据库查询延迟
	time.Sleep(100 * time.Millisecond)

	if user, exists := mockDB[id]; exists {
		return user, nil
	}
	return User{}, fmt.Errorf("用户不存在")
}

// 全局缓存实例
var userCache cache.Cache

// 获取用户信息的处理器
func getUserHandler(w http.ResponseWriter, r *http.Request) {
	// 从 URL 参数获取用户 ID
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "缺少用户ID参数", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "无效的用户ID", http.StatusBadRequest)
		return
	}

	cacheKey := fmt.Sprintf("user:%d", id)

	// 尝试从缓存获取
	if userData, exists := userCache.Get(cacheKey); exists {
		// 缓存命中
		w.Header().Set("X-Cache", "HIT")
		log.Printf("缓存命中: 用户 %d", id)
		json.NewEncoder(w).Encode(userData)
		return
	}

	// 缓存未命中，从数据库加载
	user, err := loadUserFromDB(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 存入缓存，设置 5 分钟过期时间
	err = userCache.SetWithTTL(cacheKey, user, 5*time.Minute)
	if err != nil {
		log.Printf("缓存设置失败: %v", err)
	}

	w.Header().Set("X-Cache", "MISS")
	log.Printf("缓存未命中: 从数据库加载用户 %d", id)
	json.NewEncoder(w).Encode(user)
}

// 获取缓存统计信息的处理器
func statsHandler(w http.ResponseWriter, r *http.Request) {
	stats := userCache.Stats()

	response := map[string]interface{}{
		"cache_size":  stats.Size,
		"max_size":    stats.MaxSize,
		"hits":        stats.Hits,
		"misses":      stats.Misses,
		"hit_rate":    stats.HitRate,
		"created_at":  stats.CreatedAt,
		"last_access": stats.LastAccess,
		"uptime":      time.Since(stats.CreatedAt).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 清空缓存的处理器
func clearCacheHandler(w http.ResponseWriter, r *http.Request) {
	err := userCache.Clear()
	if err != nil {
		http.Error(w, "清空缓存失败", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "缓存已清空",
		"time":    time.Now().Format(time.RFC3339),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 预热缓存 - 预先加载热门用户数据
func warmupCache() {
	log.Println("开始预热缓存...")

	hotUserIDs := []int{1, 2, 3, 4, 5} // 热门用户ID

	for _, id := range hotUserIDs {
		user, err := loadUserFromDB(id)
		if err != nil {
			log.Printf("预热失败: 用户 %d - %v", id, err)
			continue
		}

		cacheKey := fmt.Sprintf("user:%d", id)
		err = userCache.SetWithTTL(cacheKey, user, 10*time.Minute)
		if err != nil {
			log.Printf("缓存设置失败: 用户 %d - %v", id, err)
		} else {
			log.Printf("预热成功: 用户 %d", id)
		}
	}

	log.Println("缓存预热完成")
}

func main() {
	// 创建用户缓存，使用 LFU 策略
	userCache = cache.NewLFU(1000, // 最大1000个用户
		cache.WithDefaultTTL(5*time.Minute),      // 默认5分钟过期
		cache.WithCleanupInterval(2*time.Minute), // 每2分钟清理一次过期项
		cache.WithStatistics(true),               // 启用统计
	)
	defer func() {
		if err := userCache.Close(); err != nil {
			log.Printf("User cache close error: %v", err)
		}
	}()

	// 预热缓存
	warmupCache()

	// 设置路由
	http.HandleFunc("/user", getUserHandler)
	http.HandleFunc("/stats", statsHandler)
	http.HandleFunc("/clear", clearCacheHandler)

	// 健康检查
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
		})
	})

	// 首页 - 显示使用说明
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>SCache Web服务示例</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .method { color: #007bff; font-weight: bold; }
        .url { color: #333; font-family: monospace; }
        .description { color: #666; margin-top: 5px; }
    </style>
</head>
<body>
    <h1>SCache Web服务示例</h1>
    <p>这是一个展示如何在Web服务中使用SCache缓存框架的示例。</p>

    <h2>API 端点</h2>

    <div class="endpoint">
        <div><span class="method">GET</span> <span class="url">/user?id=1</span></div>
        <div class="description">获取用户信息（会自动使用缓存）</div>
    </div>

    <div class="endpoint">
        <div><span class="method">GET</span> <span class="url">/stats</span></div>
        <div class="description">获取缓存统计信息</div>
    </div>

    <div class="endpoint">
        <div><span class="method">POST</span> <span class="url">/clear</span></div>
        <div class="description">清空所有缓存</div>
    </div>

    <div class="endpoint">
        <div><span class="method">GET</span> <span class="url">/health</span></div>
        <div class="description">健康检查</div>
    </div>

    <h2>测试用户</h2>
    <p>可用的用户ID: 1, 2, 3, 4, 5</p>

    <h2>缓存信息</h2>
    <p>缓存策略: LFU (最少使用频率)</p>
    <p>最大容量: 1000 个用户</p>
    <p>默认TTL: 5 分钟</p>
    <p>清理间隔: 2 分钟</p>
</body>
</html>`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	port := ":8080"
	fmt.Printf("🚀 服务器启动在 http://localhost%s\n", port)
	fmt.Println("\n可用的API端点:")
	fmt.Println("  GET  /user?id=1        - 获取用户信息")
	fmt.Println("  GET  /stats             - 获取缓存统计")
	fmt.Println("  POST /clear             - 清空缓存")
	fmt.Println("  GET  /health            - 健康检查")
	fmt.Println("  GET  /                  - 首页说明")
	fmt.Println("\n示例:")
	fmt.Println("  curl http://localhost:8080/user?id=1")
	fmt.Println("  curl http://localhost:8080/stats")

	// 启动服务器
	log.Fatal(http.ListenAndServe(port, nil))
}

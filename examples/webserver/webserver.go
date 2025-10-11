package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"scache/cache"
)

// Response API响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// CacheStatsResponse 缓存统计响应
type CacheStatsResponse struct {
	Hits    int64    `json:"hits"`
	Misses  int64    `json:"misses"`
	Sets    int64    `json:"sets"`
	Deletes int64    `json:"deletes"`
	Size    int      `json:"size"`
	MaxSize int      `json:"max_size"`
	HitRate float64  `json:"hit_rate"`
	Keys    []string `json:"keys,omitempty"`
}

var appCache *cache.MemoryCache

func main() {
	fmt.Println("=== 缓存Web服务示例 ===")

	// 初始化应用缓存
	c := cache.NewCache(
		cache.WithMaxSize(1000),
		cache.WithDefaultExpiration(time.Minute*30),
		cache.WithCleanupInterval(time.Minute*5),
		cache.WithStats(true),
		cache.WithInitialCapacity(100),
	)
	appCache = c.(*cache.MemoryCache)

	// 设置路由
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/api/cache/set", handleSet)
	http.HandleFunc("/api/cache/get", handleGet)
	http.HandleFunc("/api/cache/delete", handleDelete)
	http.HandleFunc("/api/cache/exists", handleExists)
	http.HandleFunc("/api/cache/flush", handleFlush)
	http.HandleFunc("/api/cache/stats", handleStats)
	http.HandleFunc("/api/cache/keys", handleKeys)
	http.HandleFunc("/api/cache/size", handleSize)

	fmt.Println("\n🚀 Web服务启动在 http://localhost:8080")
	fmt.Println("\n可用的API端点:")
	fmt.Println("  GET  /                    - 主页和API说明")
	fmt.Println("  POST /api/cache/set      - 设置缓存项")
	fmt.Println("  GET  /api/cache/get      - 获取缓存项")
	fmt.Println("  DELETE /api/cache/delete - 删除缓存项")
	fmt.Println("  GET  /api/cache/exists   - 检查缓存项是否存在")
	fmt.Println("  POST /api/cache/flush    - 清空缓存")
	fmt.Println("  GET  /api/cache/stats    - 获取缓存统计")
	fmt.Println("  GET  /api/cache/keys     - 获取所有键")
	fmt.Println("  GET  /api/cache/size     - 获取缓存大小")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <title>缓存Web服务</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; }
        .endpoint { background: #f5f5f5; padding: 10px; margin: 10px 0; border-radius: 5px; }
        .method { font-weight: bold; color: #007bff; }
        pre { background: #fff; border-left: 3px solid #007bff; padding: 10px; }
    </style>
</head>
<body>
    <h1>🚀 缓存Web服务</h1>
    <p>这是一个基于Go缓存库的REST API示例服务。</p>

    <h2>📋 API端点</h2>

    <div class="endpoint">
        <span class="method">POST</span> /api/cache/set
        <pre>curl -X POST http://localhost:8080/api/cache/set \
  -H "Content-Type: application/json" \
  -d '{"key":"user1","value":"张三","ttl":"1h"}'</pre>
    </div>

    <div class="endpoint">
        <span class="method">GET</span> /api/cache/get?key=user1
        <pre>curl http://localhost:8080/api/cache/get?key=user1</pre>
    </div>

    <div class="endpoint">
        <span class="method">DELETE</span> /api/cache/delete?key=user1
        <pre>curl -X DELETE http://localhost:8080/api/cache/delete?key=user1</pre>
    </div>

    <div class="endpoint">
        <span class="method">GET</span> /api/cache/exists?key=user1
        <pre>curl http://localhost:8080/api/cache/exists?key=user1</pre>
    </div>

    <div class="endpoint">
        <span class="method">POST</span> /api/cache/flush
        <pre>curl -X POST http://localhost:8080/api/cache/flush</pre>
    </div>

    <div class="endpoint">
        <span class="method">GET</span> /api/cache/stats
        <pre>curl http://localhost:8080/api/cache/stats</pre>
    </div>

    <div class="endpoint">
        <span class="method">GET</span> /api/cache/keys
        <pre>curl http://localhost:8080/api/cache/keys</pre>
    </div>

    <div class="endpoint">
        <span class="method">GET</span> /api/cache/size
        <pre>curl http://localhost:8080/api/cache/size</pre>
    </div>

    <h2>🔧 TTL格式</h2>
    <p>支持的时间格式：</p>
    <ul>
        <li>ns (纳秒)</li>
        <li>us (微秒)</li>
        <li>ms (毫秒)</li>
        <li>s (秒)</li>
        <li>m (分钟)</li>
        <li>h (小时)</li>
    </ul>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, html)
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
		TTL   string `json:"ttl,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "无效的JSON格式", http.StatusBadRequest)
		return
	}

	if req.Key == "" {
		sendError(w, "键不能为空", http.StatusBadRequest)
		return
	}

	// 解析TTL
	var ttl time.Duration
	if req.TTL != "" {
		var err error
		ttl, err = time.ParseDuration(req.TTL)
		if err != nil {
			sendError(w, "无效的TTL格式", http.StatusBadRequest)
			return
		}
	}

	if err := appCache.Set(req.Key, req.Value, ttl); err != nil {
		sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sendResponse(w, true, "缓存设置成功", "")
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		sendError(w, "键参数是必需的", http.StatusBadRequest)
		return
	}

	value, found := appCache.Get(key)
	if !found {
		sendError(w, "缓存项未找到", http.StatusNotFound)
		return
	}

	sendResponse(w, true, value, "")
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendError(w, "只支持DELETE方法", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		sendError(w, "键参数是必需的", http.StatusBadRequest)
		return
	}

	if appCache.Delete(key) {
		sendResponse(w, true, "", "缓存删除成功")
	} else {
		sendError(w, "缓存项未找到", http.StatusNotFound)
	}
}

func handleExists(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		sendError(w, "键参数是必需的", http.StatusBadRequest)
		return
	}

	exists := appCache.Exists(key)
	sendResponse(w, true, map[string]bool{"exists": exists}, "")
}

func handleFlush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, "只支持POST方法", http.StatusMethodNotAllowed)
		return
	}

	appCache.Flush()
	sendResponse(w, true, "", "缓存已清空")
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	stats := appCache.Stats()
	response := CacheStatsResponse{
		Hits:    stats.Hits,
		Misses:  stats.Misses,
		Sets:    stats.Sets,
		Deletes: stats.Deletes,
		Size:    stats.Size,
		MaxSize: stats.MaxSize,
		HitRate: stats.HitRate,
	}

	// 添加键列表
	response.Keys = appCache.Keys()

	sendResponse(w, true, response, "")
}

func handleKeys(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	keys := appCache.Keys()
	sendResponse(w, true, keys, "")
}

func handleSize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, "只支持GET方法", http.StatusMethodNotAllowed)
		return
	}

	size := appCache.Size()
	sendResponse(w, true, map[string]int{"size": size}, "")
}

func sendResponse(w http.ResponseWriter, success bool, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")

	response := Response{
		Success: success,
		Data:    data,
	}

	if message != "" && success {
		response.Data = map[string]interface{}{"message": message, "data": data}
	}

	json.NewEncoder(w).Encode(response)
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Error:   message,
	}

	json.NewEncoder(w).Encode(response)
}

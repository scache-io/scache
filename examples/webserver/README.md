# 缓存Web服务示例

这是一个完整的缓存Web服务示例，展示了如何将scache库集成到Web应用中。

## 🚀 快速开始

### 1. 运行服务

```bash
cd examples/webserver
go run .
```

或者编译后运行：

```bash
go build -o webserver .
./webserver
```

### 2. 访问Web界面

打开浏览器访问：`http://localhost:8080`

界面功能：
- 📋 API接口文档
- 🧪 在线API测试
- 📊 实时监控面板
- 📱 响应式设计

## 🌐 远程调用

服务支持CORS跨域调用，可以从任何Web应用或客户端访问。

### JavaScript示例

```javascript
// 设置缓存
await fetch('http://localhost:8080/api/cache/set', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    key: 'user1',
    value: '张三',
    ttl: '1h'
  })
});

// 获取缓存
const response = await fetch('http://localhost:8080/api/cache/get?key=user1');
const data = await response.json();
console.log(data.data);
```

### cURL示例

```bash
# 设置缓存
curl -X POST http://localhost:8080/api/cache/set \
  -H "Content-Type: application/json" \
  -d '{"key":"test","value":"远程调用测试"}'

# 获取缓存
curl "http://localhost:8080/api/cache/get?key=test"
```

## 📁 文件结构

```
webserver/
├── main.go          # 主程序文件
├── index.html       # Web界面模板
├── API.md          # API文档
└── README.md       # 说明文件
```

## ⚙️ 配置说明

缓存配置（在main.go中）：

```go
cache.NewCache(
    cache.WithMaxSize(1000),                           // 最大1000项
    cache.WithDefaultExpiration(time.Minute*30),      // 默认30分钟过期
    cache.WithMemoryThreshold(0.7),                   // 70%内存阈值
    cache.WithBackgroundCleanup(time.Minute*2),        // 后台清理2分钟间隔
)
```

## 🔧 开发说明

### 模板系统

- 使用Go标准库`html/template`
- 支持热部署（重启生效）
- 自动查找模板文件位置

### CORS支持

- 允许所有跨域请求
- 支持预检请求(OPTIONS)
- 适合前后端分离架构

### 错误处理

- 统一JSON错误响应
- 详细的错误信息
- 合适的HTTP状态码

## 📊 API功能

| 端点 | 方法 | 功能 | 示例 |
|------|------|------|------|
| `/api/cache/set` | POST | 设置缓存 | `{"key":"k","value":"v","ttl":"1h"}` |
| `/api/cache/get` | GET | 获取缓存 | `?key=user1` |
| `/api/cache/delete` | DELETE | 删除缓存 | `?key=user1` |
| `/api/cache/exists` | GET | 检查存在 | `?key=user1` |
| `/api/cache/flush` | POST | 清空缓存 | - |
| `/api/cache/stats` | GET | 获取统计 | - |
| `/api/cache/keys` | GET | 获取所有键 | - |
| `/api/cache/size` | GET | 获取大小 | - |

详细的API文档请参考：[API.md](API.md)

## 🚦 运行端口

默认端口：`8080`

如需修改端口，编辑main.go中的端口配置：

```go
log.Fatal(http.ListenAndServe(":8080", nil))
```

## 🔍 监控功能

Web界面提供实时监控：

- 缓存大小
- 命中/未命中统计
- 命中率
- 所有缓存键列表
- 自动刷新功能

## 🛠️ 扩展功能

可以轻松扩展的功能：

1. **数据持久化** - 添加数据库存储
2. **集群支持** - 添加分布式缓存
3. **认证授权** - 添加JWT或API Key
4. **限流保护** - 添加请求限制
5. **日志记录** - 添加详细访问日志

## ❓ 常见问题

**Q: 如何从其他域名访问？**
A: 服务已启用CORS，支持所有跨域请求。

**Q: 数据会持久化吗？**
A: 不会，服务重启后数据会丢失。

**Q: 支持哪些数据类型？**
A: 当前版本支持字符串类型。

**Q: 如何修改缓存配置？**
A: 编辑main.go中的NewCache配置参数。
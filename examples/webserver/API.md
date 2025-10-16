# 缓存Web服务 API文档

## 概述

这是一个基于Go缓存库的REST API服务，支持HTTP协议的远程缓存操作。

**服务地址**: `http://localhost:8080`

**特性**:
- ✅ CORS跨域支持
- ✅ JSON格式数据交换
- ✅ 完整的CRUD操作
- ✅ 统计信息监控
- ✅ TTL过期管理

## API端点

### 1. 设置缓存项

**端点**: `POST /api/cache/set`

**请求体**:
```json
{
  "key": "string",
  "value": "string",
  "ttl": "string" // 可选，如 "30s", "5m", "1h"
}
```

**响应**:
```json
{
  "success": true,
  "message": "缓存项设置成功"
}
```

**示例**:
```bash
curl -X POST http://localhost:8080/api/cache/set \
  -H "Content-Type: application/json" \
  -d '{"key":"user1","value":"张三","ttl":"1h"}'
```

### 2. 获取缓存项

**端点**: `GET /api/cache/get?key={key}`

**响应**:
```json
{
  "success": true,
  "data": "张三",
  "message": "获取成功"
}
```

**示例**:
```bash
curl "http://localhost:8080/api/cache/get?key=user1"
```

### 3. 删除缓存项

**端点**: `DELETE /api/cache/delete?key={key}`

**响应**:
```json
{
  "success": true,
  "message": "删除成功"
}
```

**示例**:
```bash
curl -X DELETE "http://localhost:8080/api/cache/delete?key=user1"
```

### 4. 检查缓存项是否存在

**端点**: `GET /api/cache/exists?key={key}`

**响应**:
```json
{
  "success": true,
  "exists": true,
  "message": "检查完成"
}
```

**示例**:
```bash
curl "http://localhost:8080/api/cache/exists?key=user1"
```

### 5. 清空所有缓存

**端点**: `POST /api/cache/flush`

**响应**:
```json
{
  "success": true,
  "message": "缓存已清空"
}
```

**示例**:
```bash
curl -X POST http://localhost:8080/api/cache/flush
```

### 6. 获取缓存统计

**端点**: `GET /api/cache/stats`

**响应**:
```json
{
  "success": true,
  "data": {
    "hits": 150,
    "misses": 25,
    "sets": 100,
    "deletes": 10,
    "size": 90,
    "max_size": 1000,
    "hit_rate": 0.8571
  }
}
```

**示例**:
```bash
curl "http://localhost:8080/api/cache/stats"
```

### 7. 获取所有键

**端点**: `GET /api/cache/keys`

**响应**:
```json
{
  "success": true,
  "keys": ["user1", "config1", "session2"]
}
```

**示例**:
```bash
curl "http://localhost:8080/api/cache/keys"
```

### 8. 获取缓存大小

**端点**: `GET /api/cache/size`

**响应**:
```json
{
  "success": true,
  "size": 90
}
```

**示例**:
```bash
curl "http://localhost:8080/api/cache/size"
```

## TTL时间格式

支持以下时间单位：

- `ns` - 纳秒
- `us` - 微秒
- `ms` - 毫秒
- `s` - 秒
- `m` - 分钟
- `h` - 小时

**示例**:
- `30s` - 30秒
- `5m` - 5分钟
- `2h` - 2小时
- `1h30m` - 1小时30分钟

## 错误响应

所有API在出错时返回统一格式：

```json
{
  "success": false,
  "error": "错误描述信息"
}
```

常见HTTP状态码：
- `200` - 成功
- `400` - 请求参数错误
- `405` - 请求方法不允许
- `500` - 服务器内部错误

## CORS支持

服务已启用CORS跨域支持，允许来自任何域名的请求。

**允许的方法**: `GET, POST, PUT, DELETE, OPTIONS`
**允许的头部**: `Content-Type, Authorization`

## 远程调用示例

### JavaScript/Fetch API

```javascript
// 设置缓存
await fetch('http://localhost:8080/api/cache/set', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    key: 'user1',
    value: '张三',
    ttl: '1h'
  })
});

// 获取缓存
const response = await fetch('http://localhost:8080/api/cache/get?key=user1');
const data = await response.json();
console.log(data.data); // "张三"
```

### Python/requests

```python
import requests

# 设置缓存
response = requests.post('http://localhost:8080/api/cache/set', json={
    'key': 'user1',
    'value': '张三',
    'ttl': '1h'
})

# 获取缓存
response = requests.get('http://localhost:8080/api/cache/get?key=user1')
data = response.json()
print(data['data'])  # "张三"
```

### Java/OkHttp

```java
// 设置缓存
String json = "{\"key\":\"user1\",\"value\":\"张三\",\"ttl\":\"1h\"}";
RequestBody body = RequestBody.create(json, MediaType.get("application/json"));
Request request = new Request.Builder()
    .url("http://localhost:8080/api/cache/set")
    .post(body)
    .build();
try (Response response = client.newCall(request).execute()) {
    System.out.println(response.body().string());
}
```

## 在线测试

访问 `http://localhost:8080` 可以使用内置的Web界面进行在线API测试，包括：

- 🧪 API交互式测试
- 📊 实时统计监控
- 🔧 参数配置界面
- 📱 响应式设计，支持移动设备

## 性能特性

- ⚡ 高性能内存缓存
- 🔄 自动过期清理
- 📈 LRU淘汰策略
- 🛡️ 内存压力监控
- 📊 详细统计信息

## 注意事项

1. **键名限制**: 建议使用字符串键名，避免特殊字符
2. **值类型**: 当前版本仅支持字符串值
3. **TTL设置**: 不设置TTL时使用默认过期时间
4. **并发安全**: 支持高并发访问
5. **数据持久化**: 服务重启后数据会丢失
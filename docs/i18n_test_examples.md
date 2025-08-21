# HarborArk 国际化功能测试示例

## 功能说明
系统支持中文和英文两种语言，默认为中文。可以通过以下方式切换语言：

### 1. 查询参数方式
```bash
# 使用中文
curl "http://localhost:8080/api/v1/users?lang=zh-cn"

# 使用英文  
curl "http://localhost:8080/api/v1/users?lang=en-us"
```

### 2. HTTP Header方式
```bash
# 使用中文
curl -H "Accept-Language: zh-cn" "http://localhost:8080/api/v1/users"

# 使用英文
curl -H "Accept-Language: en-us" "http://localhost:8080/api/v1/users"
```

## 测试示例

### 1. 测试基础路由国际化
```bash
# 中文响应
curl "http://localhost:8080/?lang=zh-cn"
# 预期响应: {"message": "欢迎使用 HarborArk API!", ...}

# 英文响应
curl "http://localhost:8080/?lang=en-us"  
# 预期响应: {"message": "Welcome to HarborArk API!", ...}
```

### 2. 测试健康检查国际化
```bash
# 中文响应
curl "http://localhost:8080/health?lang=zh-cn"
# 预期响应: {"status": "ok", "message": "系统运行正常", ...}

# 英文响应
curl "http://localhost:8080/health?lang=en-us"
# 预期响应: {"status": "ok", "message": "System is running normally", ...}
```

### 3. 测试登录接口国际化
```bash
# 中文错误消息
curl -X POST "http://localhost:8080/api/v1/auth/login?lang=zh-cn" \
  -H "Content-Type: application/json" \
  -d '{"username": "", "password": ""}'
# 预期响应: {"code": 400, "message": "参数错误", ...}

# 英文错误消息
curl -X POST "http://localhost:8080/api/v1/auth/login?lang=en-us" \
  -H "Content-Type: application/json" \
  -d '{"username": "", "password": ""}'
# 预期响应: {"code": 400, "message": "Invalid parameters", ...}
```

### 4. 测试认证中间件国际化
```bash
# 中文错误消息 - 缺少token
curl "http://localhost:8080/api/v1/users?lang=zh-cn"
# 预期响应: {"code": 401, "message": "请求头中token为空", ...}

# 英文错误消息 - 缺少token
curl "http://localhost:8080/api/v1/users?lang=en-us"
# 预期响应: {"code": 401, "message": "Token is missing in request header", ...}
```

### 5. 测试权限中间件国际化
```bash
# 先获取普通用户token，然后尝试创建用户（需要超管权限）

# 中文权限错误
curl -X POST "http://localhost:8080/api/v1/users?lang=zh-cn" \
  -H "Authorization: Bearer NORMAL_USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "123", "group_id": 4}'
# 预期响应: {"code": 403, "message": "权限不足", ...}

# 英文权限错误
curl -X POST "http://localhost:8080/api/v1/users?lang=en-us" \
  -H "Authorization: Bearer NORMAL_USER_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "123", "group_id": 4}'
# 预期响应: {"code": 403, "message": "Permission denied", ...}
```

## 支持的语言代码
- `zh-cn`, `zh`, `chinese` - 中文
- `en-us`, `en`, `english` - 英文
- 默认语言：中文 (`zh-cn`)

## 翻译键值对照表

| 键名 | 中文 | 英文 |
|------|------|------|
| welcome | 欢迎使用 HarborArk API! | Welcome to HarborArk API! |
| health_ok | 系统运行正常 | System is running normally |
| login_success | 登录成功 | Login successful |
| invalid_params | 参数错误 | Invalid parameters |
| invalid_token | 无效的token | Invalid token |
| token_missing | 请求头中token为空 | Token is missing in request header |
| permission_denied | 权限不足 | Permission denied |
| user_not_found | 用户不存在 | User not found |
| create_success | 创建成功 | Created successfully |
| ... | ... | ... |

完整的翻译文件位于：
- `translations/zh-cn.json` - 中文翻译
- `translations/en-us.json` - 英文翻译
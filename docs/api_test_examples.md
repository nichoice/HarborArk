# HarborArk 用户管理系统 API 测试示例

## 基本信息
- 服务器地址：http://localhost:8080
- API文档：http://localhost:8080/swagger/index.html
- 默认管理员：用户名 `admin`，密码 `admin123`

## API 测试示例

### 1. 用户登录
```bash
curl -X POST "http://localhost:8080/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### 2. 获取用户组列表
```bash
# 先登录获取token，然后使用token访问
curl -X GET "http://localhost:8080/api/v1/user-groups" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json"
```

### 3. 获取用户列表
```bash
curl -X GET "http://localhost:8080/api/v1/users" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json"
```

### 4. 创建新用户（需要超级管理员权限）
```bash
curl -X POST "http://localhost:8080/api/v1/users" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "test123",
    "group_id": 4
  }'
```

### 5. 更新用户信息
```bash
curl -X PUT "http://localhost:8080/api/v1/users/2" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -H "Content-Type: application/json" \
  -d '{
    "is_enabled": false
  }'
```

## 用户组级别说明
- 1: 超级管理员 - 拥有所有权限
- 2: 运维管理员 - 运维管理权限
- 3: 审计管理员 - 审计管理权限  
- 4: 普通用户 - 基础权限，会同步创建Linux用户（nologin类型）

## 权限说明
- 只有超级管理员（level=1）可以进行用户和用户组的增删改操作
- 所有已认证用户都可以查看用户和用户组列表
- JWT token有效期为24小时

## 特殊功能
- 普通用户（group_id=4）创建时会自动在Linux系统中创建对应用户（nologin类型）
- 用户启用/禁用会同步更新Linux用户的shell状态
- 删除普通用户时会同步删除Linux系统用户
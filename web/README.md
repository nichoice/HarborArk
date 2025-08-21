# HarborArk 前端管理系统

基于 Vue3 + TDesign + TypeScript 构建的用户管理系统前端。

## 功能特性

- 🔐 **用户认证**: JWT登录认证
- 👥 **用户管理**: 用户的增删改查操作
- 🏷️ **用户组管理**: 用户组的增删改查操作
- 🌐 **国际化支持**: 中英文切换
- 📱 **响应式设计**: 适配各种屏幕尺寸
- 🎨 **TDesign UI**: 腾讯TDesign设计语言

## 技术栈

- **框架**: Vue 3 + TypeScript
- **构建工具**: Vite
- **UI组件库**: TDesign Vue Next
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios

## 快速开始

### 1. 安装依赖

```bash
cd web
npm install
```

### 2. 启动开发服务器

```bash
npm run dev
```

前端服务将在 http://localhost:3000 启动

### 3. 构建生产版本

```bash
npm run build
```

## 项目结构

```
web/
├── src/
│   ├── api/              # API接口
│   ├── components/       # 公共组件
│   ├── layouts/          # 布局组件
│   ├── router/           # 路由配置
│   ├── stores/           # 状态管理
│   ├── types/            # TypeScript类型定义
│   ├── views/            # 页面组件
│   ├── App.vue           # 根组件
│   └── main.ts           # 入口文件
├── public/               # 静态资源
├── index.html            # HTML模板
├── package.json          # 项目配置
├── tsconfig.json         # TypeScript配置
└── vite.config.ts        # Vite配置
```

## 页面说明

### 登录页面 (`/login`)
- 用户名/密码登录
- JWT token认证
- 自动跳转到仪表盘

### 仪表盘 (`/dashboard`)
- 系统概览
- 统计数据展示
- 快速操作入口

### 用户管理 (`/system/users`)
- 用户列表查看
- 新增/编辑/删除用户
- 用户状态管理
- 用户组分配

### 用户组管理 (`/system/user-groups`)
- 用户组列表查看
- 新增/编辑/删除用户组
- 权限配置管理

## API接口

前端通过代理访问后端API：
- 开发环境: `http://localhost:3000/api` -> `http://localhost:8080/api`
- 生产环境: 需要配置nginx反向代理

## 默认登录账号

- 用户名: `admin`
- 密码: `admin123`

## 开发说明

1. 确保后端服务已启动 (端口8080)
2. 前端开发服务器会自动代理API请求到后端
3. 支持热重载，修改代码后自动刷新页面
4. 使用TDesign组件库，遵循设计规范

## 构建部署

1. 执行构建命令生成dist目录
2. 将dist目录部署到Web服务器
3. 配置nginx反向代理API请求到后端服务
4. 确保前端路由支持History模式
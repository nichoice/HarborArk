#!/bin/bash

# HarborArk 系统启动脚本

echo "🚀 启动 HarborArk 管理系统..."

# 检查Go后端是否已编译
if [ ! -f "./harborark" ]; then
    echo "📦 编译Go后端..."
    go build -o harborark main.go
fi

# 检查前端依赖是否已安装
if [ ! -d "./web/node_modules" ]; then
    echo "📦 安装前端依赖..."
    cd web && npm install && cd ..
fi

echo "🔧 启动后端服务 (端口: 8080)..."
./harborark server &
BACKEND_PID=$!

# 等待后端启动
sleep 3

echo "🎨 启动前端服务 (端口: 3000)..."
cd web && npm run dev &
FRONTEND_PID=$!

echo ""
echo "✅ 系统启动完成!"
echo ""
echo "📋 访问地址:"
echo "   前端管理系统: http://localhost:3000"
echo "   后端API文档:  http://localhost:8080/swagger/index.html"
echo "   健康检查:     http://localhost:8080/health"
echo ""
echo "🔑 默认登录账号:"
echo "   用户名: admin"
echo "   密码:   admin123"
echo ""
echo "⚠️  按 Ctrl+C 停止所有服务"

# 等待用户中断
trap "echo '🛑 正在停止服务...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null; exit" INT

# 保持脚本运行
wait
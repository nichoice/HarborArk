#!/bin/bash

BASE_URL="http://localhost:8080"

echo "=== HarborArk 国际化功能测试 ==="
echo

# 等待服务器启动
sleep 2

echo "1. 测试基础路由国际化"
echo "中文响应:"
curl -s "$BASE_URL/?lang=zh-cn" | jq '.'
echo
echo "英文响应:"
curl -s "$BASE_URL/?lang=en-us" | jq '.'
echo

echo "2. 测试健康检查国际化"
echo "中文响应:"
curl -s "$BASE_URL/health?lang=zh-cn" | jq '.'
echo
echo "英文响应:"
curl -s "$BASE_URL/health?lang=en-us" | jq '.'
echo

echo "3. 测试登录接口参数错误国际化"
echo "中文错误消息:"
curl -s -X POST "$BASE_URL/api/v1/auth/login?lang=zh-cn" \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.'
echo
echo "英文错误消息:"
curl -s -X POST "$BASE_URL/api/v1/auth/login?lang=en-us" \
  -H "Content-Type: application/json" \
  -d '{}' | jq '.'
echo

echo "4. 测试认证中间件国际化（缺少token）"
echo "中文错误消息:"
curl -s "$BASE_URL/api/v1/users?lang=zh-cn" | jq '.'
echo
echo "英文错误消息:"
curl -s "$BASE_URL/api/v1/users?lang=en-us" | jq '.'
echo

echo "5. 测试Header方式语言切换"
echo "使用Accept-Language头 - 中文:"
curl -s -H "Accept-Language: zh-cn" "$BASE_URL/" | jq '.'
echo
echo "使用Accept-Language头 - 英文:"
curl -s -H "Accept-Language: en-us" "$BASE_URL/" | jq '.'
echo

echo "=== 国际化功能测试完成 ==="
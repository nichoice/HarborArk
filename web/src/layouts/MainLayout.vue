<template>
  <t-layout class="main-layout">
    <t-aside class="sidebar" :width="240">
      <div class="sidebar-header">
        <div class="logo">
          <t-icon name="cloud" size="24px" />
          <span class="logo-text">HarborArk</span>
        </div>
      </div>
      
      <t-menu 
        v-model:value="activeMenu" 
        :default-expanded="['system-management', 'file-management']"
        theme="light"
        @change="handleMenuChange"
      >
        <t-submenu value="file-management" title="文件管理">
          <template #icon>
            <t-icon name="folder" />
          </template>
          <t-menu-item value="file-browser">
            <template #icon>
              <t-icon name="browse" />
            </template>
            文件浏览器
          </t-menu-item>
          <t-menu-item value="file-upload">
            <template #icon>
              <t-icon name="upload" />
            </template>
            文件上传
          </t-menu-item>
        </t-submenu>
        
        <t-submenu value="system-management" title="系统管理">
          <template #icon>
            <t-icon name="setting" />
          </template>
          <t-menu-item value="users">
            <template #icon>
              <t-icon name="user-list" />
            </template>
            用户管理
          </t-menu-item>
          <t-menu-item value="user-groups">
            <template #icon>
              <t-icon name="usergroup" />
            </template>
            用户组管理
          </t-menu-item>
          <t-menu-item value="audit-logs">
            <template #icon>
              <t-icon name="file-text" />
            </template>
            审计日志
          </t-menu-item>
          <t-menu-item value="settings">
            <template #icon>
              <t-icon name="tools" />
            </template>
            系统设置
          </t-menu-item>
        </t-submenu>
      </t-menu>
    </t-aside>
    
    <t-layout>
      <t-header class="header">
        <div class="header-content">
          <div class="header-left">
            <t-breadcrumb>
              <t-breadcrumb-item>{{ getCurrentPageTitle() }}</t-breadcrumb-item>
            </t-breadcrumb>
          </div>
          <div class="header-actions">
            <t-dropdown :options="userMenuOptions" @click="handleUserMenu">
              <t-button variant="text" class="user-btn">
                <t-icon name="user-circle" size="20px" />
                <span>管理员</span>
                <t-icon name="chevron-down" size="16px" />
              </t-button>
            </t-dropdown>
          </div>
        </div>
      </t-header>
      
      <t-content class="content">
        <div class="content-wrapper">
          <router-view />
        </div>
      </t-content>
    </t-layout>
  </t-layout>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const activeMenu = ref('users')

const userMenuOptions = [
  { content: '个人设置', value: 'profile' },
  { content: '退出登录', value: 'logout' }
]

const handleMenuChange = (value: string) => {
  console.log('菜单变化:', value)
  if (value) {
    router.push(`/${value}`)
  }
}

const handleMenuClick = (value: any) => {
  console.log('菜单点击:', value)
  // TDesign 菜单点击事件可能传递的是对象
  const menuValue = typeof value === 'string' ? value : value.value || value
  router.push(`/${menuValue}`)
}

const handleUserMenu = (data: any) => {
  if (data.value === 'logout') {
    // 处理退出登录
    console.log('退出登录')
  }
}

// 获取当前页面标题
const getCurrentPageTitle = () => {
  const titleMap: Record<string, string> = {
    '/users': '用户管理',
    '/user-groups': '用户组管理',
    '/file-browser': '文件浏览器',
    '/file-upload': '文件上传',
    '/audit-logs': '审计日志',
    '/settings': '系统设置'
  }
  return titleMap[route.path] || '首页'
}

// 监听路由变化更新菜单状态
watch(() => route.path, (newPath) => {
  const pathMap: Record<string, string> = {
    '/users': 'users',
    '/user-groups': 'user-groups',
    '/file-browser': 'file-browser',
    '/file-upload': 'file-upload',
    '/audit-logs': 'audit-logs',
    '/settings': 'settings'
  }
  activeMenu.value = pathMap[newPath] || 'users'
}, { immediate: true })
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.sidebar {
  background: #fff;
  border-right: 1px solid #e7e7e7;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 16px 20px;
  border-bottom: 1px solid #e7e7e7;
  background: #fafafa;
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 18px;
  font-weight: 600;
  color: #0052d9;
}

.logo-text {
  color: #262626;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e7e7e7;
  padding: 0;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 56px;
  padding: 0 24px;
}

.header-left {
  flex: 1;
}

.header-actions {
  display: flex;
  align-items: center;
}

.user-btn {
  display: flex;
  align-items: center;
  gap: 4px;
}

.content {
  background: #f5f5f5;
  padding: 0;
}

.content-wrapper {
  padding: 24px;
  min-height: calc(100vh - 56px);
}

:deep(.t-menu) {
  border-right: none;
}

:deep(.t-menu .t-menu__item) {
  margin: 2px 8px;
  border-radius: 6px;
}

:deep(.t-menu .t-submenu__title) {
  margin: 2px 8px;
  border-radius: 6px;
}
</style>

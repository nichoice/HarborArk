<template>
  <t-layout class="main-layout">
    <!-- 侧边栏 -->
    <t-aside width="240px" class="sidebar">
      <div class="logo">
        <h2>HarborArk</h2>
      </div>
      
      <t-menu
        v-model:value="activeMenu"
        :default-expanded="['system']"
        @change="handleMenuChange"
        class="sidebar-menu"
      >
        <t-menu-item value="dashboard">
          <template #icon>
            <t-icon name="dashboard" />
          </template>
          仪表盘
        </t-menu-item>
        
        <t-submenu value="system">
          <template #icon>
            <t-icon name="setting" />
          </template>
          <template #title>系统管理</template>
          
          <t-menu-item value="users">
            <template #icon>
              <t-icon name="user" />
            </template>
            用户管理
          </t-menu-item>
          
          <t-menu-item value="user-groups">
            <template #icon>
              <t-icon name="usergroup" />
            </template>
            用户组管理
          </t-menu-item>
        </t-submenu>
      </t-menu>
    </t-aside>

    <!-- 主内容区 -->
    <t-layout>
      <!-- 顶部导航 -->
      <t-header class="header">
        <div class="header-content">
          <div class="breadcrumb">
            <t-breadcrumb>
              <t-breadcrumb-item v-for="item in breadcrumbs" :key="item.path">
                {{ item.title }}
              </t-breadcrumb-item>
            </t-breadcrumb>
          </div>
          
          <div class="user-info">
            <t-dropdown>
              <div class="user-avatar">
                <t-avatar size="small">{{ userInitial }}</t-avatar>
                <span class="username">{{ authStore.user?.full_name || authStore.user?.username }}</span>
                <t-icon name="chevron-down" />
              </div>
              
              <t-dropdown-menu>
                <t-dropdown-item @click="handleLogout">
                  <t-icon name="logout" />
                  退出登录
                </t-dropdown-item>
              </t-dropdown-menu>
            </t-dropdown>
          </div>
        </div>
      </t-header>

      <!-- 内容区域 -->
      <t-content class="content">
        <router-view />
      </t-content>
    </t-layout>
  </t-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const activeMenu = ref('dashboard')

// 用户名首字母
const userInitial = computed(() => {
  const name = authStore.user?.full_name || authStore.user?.username || 'U'
  return name.charAt(0).toUpperCase()
})

// 面包屑导航
const breadcrumbs = computed(() => {
  const crumbs = []
  
  if (route.path.startsWith('/system')) {
    crumbs.push({ title: '系统管理', path: '/system' })
    
    if (route.path === '/system/users') {
      crumbs.push({ title: '用户管理', path: '/system/users' })
    } else if (route.path === '/system/user-groups') {
      crumbs.push({ title: '用户组管理', path: '/system/user-groups' })
    }
  } else if (route.path === '/dashboard') {
    crumbs.push({ title: '仪表盘', path: '/dashboard' })
  }
  
  return crumbs
})

// 监听路由变化更新菜单
watch(() => route.path, (newPath) => {
  if (newPath === '/dashboard') {
    activeMenu.value = 'dashboard'
  } else if (newPath === '/system/users') {
    activeMenu.value = 'users'
  } else if (newPath === '/system/user-groups') {
    activeMenu.value = 'user-groups'
  }
}, { immediate: true })

// 菜单点击处理
const handleMenuChange = (value: string) => {
  switch (value) {
    case 'dashboard':
      router.push('/dashboard')
      break
    case 'users':
      router.push('/system/users')
      break
    case 'user-groups':
      router.push('/system/user-groups')
      break
  }
}

// 退出登录
const handleLogout = () => {
  authStore.logout()
  MessagePlugin.success('已退出登录')
  router.push('/login')
}
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.sidebar {
  background: #f8f9fa;
  color: #333;
  border-right: 1px solid #e9ecef;
}

.logo {
  padding: 16px 24px;
  text-align: center;
  border-bottom: 1px solid #e9ecef;
}

.logo h2 {
  color: #333;
  margin: 0;
  font-size: 18px;
  font-weight: 600;
}

.sidebar-menu {
  background: transparent;
  border: none;
}

.sidebar-menu :deep(.t-menu-item),
.sidebar-menu :deep(.t-submenu__title) {
  color: #333 !important;
}

.sidebar-menu :deep(.t-menu-item:hover),
.sidebar-menu :deep(.t-submenu__title:hover) {
  background: #e9ecef !important;
  color: #333 !important;
}

.sidebar-menu :deep(.t-menu-item.t-is-active) {
  background: #1890ff !important;
  color: white !important;
}

.sidebar-menu :deep(.t-submenu .t-menu-item) {
  color: #666 !important;
}

.sidebar-menu :deep(.t-submenu .t-menu-item:hover) {
  background: #e9ecef !important;
  color: #333 !important;
}

.sidebar-menu :deep(.t-submenu .t-menu-item.t-is-active) {
  background: #1890ff !important;
  color: white !important;
}

.sidebar-menu :deep(.t-icon) {
  color: inherit !important;
}

.header {
  background: white;
  border-bottom: 1px solid #f0f0f0;
  padding: 0 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 100%;
}

.user-info {
  display: flex;
  align-items: center;
}

.user-avatar {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 8px 12px;
  border-radius: 6px;
  transition: background-color 0.2s;
}

.user-avatar:hover {
  background: #f5f5f5;
}

.username {
  font-size: 14px;
  color: #333;
}

.content {
  padding: 24px;
  background: #f5f5f5;
  overflow-y: auto;
}
</style>
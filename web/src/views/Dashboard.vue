<template>
  <div class="dashboard">
    <div class="welcome-card">
      <t-card>
        <h2>欢迎使用 HarborArk 管理系统</h2>
        <p>这是一个基于 Go + Vue3 + TDesign 构建的用户管理系统</p>
        
        <div class="stats-grid">
          <div class="stat-item">
            <t-icon name="user" class="stat-icon" />
            <div class="stat-content">
              <div class="stat-number">{{ stats.userCount }}</div>
              <div class="stat-label">用户总数</div>
            </div>
          </div>
          
          <div class="stat-item">
            <t-icon name="usergroup" class="stat-icon" />
            <div class="stat-content">
              <div class="stat-number">{{ stats.groupCount }}</div>
              <div class="stat-label">用户组数</div>
            </div>
          </div>
          
          <div class="stat-item">
            <t-icon name="time" class="stat-icon" />
            <div class="stat-content">
              <div class="stat-number">{{ stats.onlineCount }}</div>
              <div class="stat-label">在线用户</div>
            </div>
          </div>
        </div>
      </t-card>
    </div>

    <div class="quick-actions">
      <t-card title="快速操作">
        <div class="action-buttons">
          <t-button theme="primary" @click="$router.push('/system/users')">
            <t-icon name="user" />
            用户管理
          </t-button>
          
          <t-button theme="default" @click="$router.push('/system/user-groups')">
            <t-icon name="usergroup" />
            用户组管理
          </t-button>
        </div>
      </t-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const stats = ref({
  userCount: 0,
  groupCount: 0,
  onlineCount: 1
})

onMounted(() => {
  // 这里可以调用API获取统计数据
  // 暂时使用模拟数据
  stats.value = {
    userCount: 12,
    groupCount: 3,
    onlineCount: 1
  }
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
}

.welcome-card {
  margin-bottom: 24px;
}

.welcome-card h2 {
  margin: 0 0 8px 0;
  color: #1f2937;
}

.welcome-card p {
  margin: 0 0 24px 0;
  color: #6b7280;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 24px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  background: #f8fafc;
  border-radius: 8px;
}

.stat-icon {
  font-size: 32px;
  color: #1890ff;
}

.stat-content {
  flex: 1;
}

.stat-number {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: #6b7280;
}

.action-buttons {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.action-buttons .t-button {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
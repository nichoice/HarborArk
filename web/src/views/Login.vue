<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h1>HarborArk 管理系统</h1>
        <p>欢迎登录</p>
      </div>
      
      <t-form
        ref="formRef"
        :data="formData"
        :rules="rules"
        @submit="handleSubmit"
        class="login-form"
      >
        <t-form-item name="username">
          <t-input
            v-model="formData.username"
            placeholder="请输入用户名"
            size="large"
            clearable
          >
            <template #prefix-icon>
              <t-icon name="user" />
            </template>
          </t-input>
        </t-form-item>
        
        <t-form-item name="password">
          <t-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            clearable
          >
            <template #prefix-icon>
              <t-icon name="lock-on" />
            </template>
          </t-input>
        </t-form-item>
        
        <t-form-item>
          <t-button
            type="submit"
            theme="primary"
            size="large"
            :loading="loading"
            block
          >
            登录
          </t-button>
        </t-form-item>
      </t-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest } from '@/types/auth'

const router = useRouter()
const authStore = useAuthStore()

const formRef = ref()
const loading = ref(false)

const formData = reactive<LoginRequest>({
  username: '',
  password: ''
})

const rules = {
  username: [
    { required: true, message: '请输入用户名', type: 'error' }
  ],
  password: [
    { required: true, message: '请输入密码', type: 'error' }
  ]
}

const handleSubmit = async ({ validateResult }: any) => {
  if (validateResult === true) {
    loading.value = true
    try {
      await authStore.login(formData)
      MessagePlugin.success('登录成功')
      router.push('/')
    } catch (error: any) {
      console.error('登录失败:', error)
    } finally {
      loading.value = false
    }
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  padding: 40px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-header h1 {
  font-size: 24px;
  font-weight: 600;
  color: #1f2937;
  margin-bottom: 8px;
}

.login-header p {
  color: #6b7280;
  font-size: 14px;
}

.login-form {
  width: 100%;
}

.login-form :deep(.t-form-item) {
  margin-bottom: 24px;
}

.login-form :deep(.t-form-item:last-child) {
  margin-bottom: 0;
}
</style>
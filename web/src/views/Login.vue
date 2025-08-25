<template>
  <div class="login-page">
    <div class="login-container">
      <div class="login-card">
        <div class="login-header">
          <div class="logo">
            <t-icon name="cloud" size="32px" />
            <h1>HarborArk</h1>
          </div>
          <p class="subtitle">安全可靠的文件管理系统</p>
        </div>
        
        <t-form
          ref="formRef"
          :model="loginForm"
          :rules="formRules"
          label-width="0"
          @submit="handleLogin"
        >
          <t-form-item name="username">
            <t-input
              v-model="loginForm.username"
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
              v-model="loginForm.password"
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
            <t-checkbox v-model="rememberMe">记住我</t-checkbox>
          </t-form-item>
          
          <t-form-item>
            <t-button
              theme="primary"
              type="submit"
              size="large"
              :loading="loading"
              block
            >
              登录
            </t-button>
          </t-form-item>
        </t-form>
        
        <div v-if="errorMessage" class="error-message">
          <t-alert theme="error" :message="errorMessage" />
        </div>
      </div>
      
      <div class="login-footer">
        <p>&copy; 2024 HarborArk. All rights reserved.</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { MessagePlugin } from 'tdesign-vue-next'

const router = useRouter()
const formRef = ref()
const loading = ref(false)
const errorMessage = ref('')
const rememberMe = ref(false)

// 登录表单
const loginForm = reactive({
  username: 'admin',
  password: 'admin123'
})

// 表单验证规则
const formRules = {
  username: [
    { required: true, message: '请输入用户名', type: 'error' }
  ],
  password: [
    { required: true, message: '请输入密码', type: 'error' }
  ]
}

// 处理登录
const handleLogin = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return
  
  loading.value = true
  errorMessage.value = ''
  
  try {
    const response = await fetch('/api/v1/auth/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        username: loginForm.username,
        password: loginForm.password
      })
    })
    
    const data = await response.json()
    
    if (!response.ok || data?.code !== 200) {
      throw new Error(data?.message || '登录失败')
    }
    
    const token = data?.data?.token || data?.token || ''
    localStorage.setItem('token', token)
    
    if (rememberMe.value) {
      localStorage.setItem('remember_username', loginForm.username)
    } else {
      localStorage.removeItem('remember_username')
    }
    
    MessagePlugin.success('登录成功')
    router.push('/')
    
  } catch (error: any) {
    errorMessage.value = error.message || '登录失败，请检查用户名和密码'
  } finally {
    loading.value = false
  }
}

// 组件挂载时检查是否记住用户名
const savedUsername = localStorage.getItem('remember_username')
if (savedUsername) {
  loginForm.username = savedUsername
  rememberMe.value = true
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.login-container {
  width: 100%;
  max-width: 400px;
}

.login-card {
  background: #fff;
  border-radius: 12px;
  padding: 40px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  margin-bottom: 20px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.logo {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin-bottom: 8px;
}

.logo h1 {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: #262626;
}

.subtitle {
  margin: 0;
  color: #8a8e99;
  font-size: 14px;
}

.error-message {
  margin-top: 16px;
}

.login-footer {
  text-align: center;
}

.login-footer p {
  margin: 0;
  color: rgba(255, 255, 255, 0.8);
  font-size: 14px;
}

:deep(.t-form-item) {
  margin-bottom: 20px;
}

:deep(.t-form-item:last-child) {
  margin-bottom: 0;
}
</style>
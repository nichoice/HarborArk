<template>
  <div class="users-page">
    <t-card>
      <template #header>
        <div class="page-header">
          <h3>用户管理</h3>
          <t-button theme="primary" @click="showCreateDialog = true">
            <t-icon name="add" />
            新增用户
          </t-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <t-input
          v-model="searchKeyword"
          placeholder="搜索用户名或邮箱"
          clearable
          @enter="handleSearch"
          style="width: 300px"
        >
          <template #suffix-icon>
            <t-icon name="search" @click="handleSearch" />
          </template>
        </t-input>
        <t-button @click="handleSearch">搜索</t-button>
        <t-button theme="default" @click="handleReset">重置</t-button>
      </div>

      <!-- 用户表格 -->
      <t-table
        :data="users"
        :columns="columns"
        :loading="loading"
        row-key="id"
        :pagination="pagination"
        @page-change="handlePageChange"
      >
        <template #is_active="{ row }">
          <t-tag :theme="row.is_active ? 'success' : 'danger'">
            {{ row.is_active ? '启用' : '禁用' }}
          </t-tag>
        </template>

        <template #user_group="{ row }">
          <t-tag theme="primary">{{ row.user_group?.name || '未分组' }}</t-tag>
        </template>

        <template #created_at="{ row }">
          {{ formatDate(row.created_at) }}
        </template>

        <template #operation="{ row }">
          <t-space>
            <t-button theme="primary" variant="text" @click="handleEdit(row)">
              编辑
            </t-button>
            <t-button theme="danger" variant="text" @click="handleDelete(row)">
              删除
            </t-button>
          </t-space>
        </template>
      </t-table>
    </t-card>

    <!-- 新增/编辑用户对话框 -->
    <t-dialog
      v-model:visible="showCreateDialog"
      :header="editingUser ? '编辑用户' : '新增用户'"
      width="600px"
      @confirm="handleSubmit"
      @cancel="handleCancel"
    >
      <t-form
        ref="formRef"
        :data="formData"
        :rules="formRules"
        label-width="100px"
      >
        <t-form-item label="用户名" name="username">
          <t-input
            v-model="formData.username"
            placeholder="请输入用户名"
            :disabled="!!editingUser"
          />
        </t-form-item>

        <t-form-item label="密码" name="password" v-if="!editingUser">
          <t-input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码"
          />
        </t-form-item>

        <t-form-item label="邮箱" name="email">
          <t-input
            v-model="formData.email"
            placeholder="请输入邮箱"
          />
        </t-form-item>

        <t-form-item label="姓名" name="full_name">
          <t-input
            v-model="formData.full_name"
            placeholder="请输入姓名"
          />
        </t-form-item>

        <t-form-item label="用户组" name="user_group_id">
          <t-select
            v-model="formData.user_group_id"
            placeholder="请选择用户组"
            clearable
          >
            <t-option
              v-for="group in userGroups"
              :key="group.id"
              :value="group.id"
              :label="group.name"
            />
          </t-select>
        </t-form-item>

        <t-form-item label="状态" name="is_active" v-if="editingUser">
          <t-switch v-model="formData.is_active" />
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import { usersApi } from '@/api/users'
import { userGroupsApi } from '@/api/userGroups'
import type { User, UserGroup } from '@/types/auth'

const loading = ref(false)
const users = ref<User[]>([])
const userGroups = ref<UserGroup[]>([])
const searchKeyword = ref('')
const showCreateDialog = ref(false)
const editingUser = ref<User | null>(null)
const formRef = ref()

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  username: '',
  password: '',
  email: '',
  full_name: '',
  user_group_id: undefined as number | undefined,
  is_active: true
})

const formRules = {
  username: [
    { required: true, message: '请输入用户名', type: 'error' }
  ],
  password: [
    { required: true, message: '请输入密码', type: 'error' }
  ],
  email: [
    { required: true, message: '请输入邮箱', type: 'error' },
    { email: true, message: '请输入正确的邮箱格式', type: 'error' }
  ],
  full_name: [
    { required: true, message: '请输入姓名', type: 'error' }
  ]
}

const columns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'username', title: '用户名', width: 120 },
  { colKey: 'email', title: '邮箱', width: 200 },
  { colKey: 'full_name', title: '姓名', width: 120 },
  { colKey: 'user_group', title: '用户组', width: 120 },
  { colKey: 'is_active', title: '状态', width: 100 },
  { colKey: 'created_at', title: '创建时间', width: 160 },
  { colKey: 'operation', title: '操作', width: 150, fixed: 'right' }
]

// 获取用户列表
const fetchUsers = async () => {
  loading.value = true
  try {
    const response = await usersApi.getUsers({
      page: pagination.current,
      page_size: pagination.pageSize
    })
    users.value = response.data?.list || []
    pagination.total = response.data?.total || 0
  } catch (error) {
    console.error('获取用户列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取用户组列表
const fetchUserGroups = async () => {
  try {
    const response = await userGroupsApi.getUserGroups()
    userGroups.value = response.data || []
  } catch (error) {
    console.error('获取用户组列表失败:', error)
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 搜索
const handleSearch = () => {
  pagination.current = 1
  fetchUsers()
}

// 重置
const handleReset = () => {
  searchKeyword.value = ''
  pagination.current = 1
  fetchUsers()
}

// 分页变化
const handlePageChange = (pageInfo: any) => {
  pagination.current = pageInfo.current
  pagination.pageSize = pageInfo.pageSize
  fetchUsers()
}

// 编辑用户
const handleEdit = (user: User) => {
  editingUser.value = user
  Object.assign(formData, {
    username: user.username,
    email: user.email,
    full_name: user.full_name,
    user_group_id: user.user_group_id,
    is_active: user.is_active,
    password: ''
  })
  showCreateDialog.value = true
}

// 删除用户
const handleDelete = (user: User) => {
  const dialog = DialogPlugin.confirm({
    header: '确认删除',
    body: `确定要删除用户 "${user.username}" 吗？此操作不可恢复。`,
    onConfirm: async () => {
      try {
        await usersApi.deleteUser(user.id)
        MessagePlugin.success('删除成功')
        fetchUsers()
        dialog.destroy()
      } catch (error) {
        console.error('删除用户失败:', error)
      }
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  try {
    if (editingUser.value) {
      // 编辑用户
      await usersApi.updateUser(editingUser.value.id, {
        email: formData.email,
        full_name: formData.full_name,
        user_group_id: formData.user_group_id,
        is_active: formData.is_active
      })
      MessagePlugin.success('更新成功')
    } else {
      // 新增用户
      await usersApi.createUser({
        username: formData.username,
        password: formData.password,
        email: formData.email,
        full_name: formData.full_name,
        user_group_id: formData.user_group_id!
      })
      MessagePlugin.success('创建成功')
    }
    
    showCreateDialog.value = false
    fetchUsers()
  } catch (error) {
    console.error('操作失败:', error)
  }
}

// 取消操作
const handleCancel = () => {
  showCreateDialog.value = false
  editingUser.value = null
  Object.assign(formData, {
    username: '',
    password: '',
    email: '',
    full_name: '',
    user_group_id: undefined,
    is_active: true
  })
}

onMounted(() => {
  fetchUsers()
  fetchUserGroups()
})
</script>

<style scoped>
.users-page {
  height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-header h3 {
  margin: 0;
}

.search-bar {
  display: flex;
  gap: 16px;
  margin-bottom: 16px;
  align-items: center;
}
</style>
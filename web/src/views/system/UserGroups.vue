<template>
  <div class="user-groups-page">
    <t-card>
      <template #header>
        <div class="page-header">
          <h3>用户组管理</h3>
          <t-button theme="primary" @click="showCreateDialog = true">
            <t-icon name="add" />
            新增用户组
          </t-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <div class="search-bar">
        <t-input
          v-model="searchKeyword"
          placeholder="搜索用户组名称"
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

      <!-- 用户组表格 -->
      <t-table
        :data="userGroups"
        :columns="columns"
        :loading="loading"
        row-key="id"
        :pagination="pagination"
        @page-change="handlePageChange"
      >
        <template #permissions="{ row }">
          <t-space>
            <t-tag
              v-for="permission in row.permissions"
              :key="permission"
              theme="primary"
              variant="light"
            >
              {{ getPermissionLabel(permission) }}
            </t-tag>
          </t-space>
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

    <!-- 新增/编辑用户组对话框 -->
    <t-dialog
      v-model:visible="showCreateDialog"
      :header="editingUserGroup ? '编辑用户组' : '新增用户组'"
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
        <t-form-item label="组名" name="name">
          <t-input
            v-model="formData.name"
            placeholder="请输入用户组名称"
          />
        </t-form-item>

        <t-form-item label="描述" name="description">
          <t-textarea
            v-model="formData.description"
            placeholder="请输入用户组描述"
            :autosize="{ minRows: 3, maxRows: 6 }"
          />
        </t-form-item>

        <t-form-item label="权限" name="permissions">
          <t-checkbox-group v-model="formData.permissions">
            <div class="permissions-grid">
              <t-checkbox
                v-for="permission in availablePermissions"
                :key="permission.value"
                :value="permission.value"
              >
                {{ permission.label }}
              </t-checkbox>
            </div>
          </t-checkbox-group>
        </t-form-item>
      </t-form>
    </t-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { MessagePlugin, DialogPlugin } from 'tdesign-vue-next'
import { userGroupsApi } from '@/api/userGroups'
import type { UserGroup } from '@/types/auth'

const loading = ref(false)
const userGroups = ref<UserGroup[]>([])
const searchKeyword = ref('')
const showCreateDialog = ref(false)
const editingUserGroup = ref<UserGroup | null>(null)
const formRef = ref()

const pagination = reactive({
  current: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive({
  name: '',
  description: '',
  permissions: [] as string[]
})

const formRules = {
  name: [
    { required: true, message: '请输入用户组名称', type: 'error' }
  ],
  description: [
    { required: true, message: '请输入用户组描述', type: 'error' }
  ]
}

const availablePermissions = [
  { value: 'user:read', label: '查看用户' },
  { value: 'user:create', label: '创建用户' },
  { value: 'user:update', label: '更新用户' },
  { value: 'user:delete', label: '删除用户' },
  { value: 'group:read', label: '查看用户组' },
  { value: 'group:create', label: '创建用户组' },
  { value: 'group:update', label: '更新用户组' },
  { value: 'group:delete', label: '删除用户组' },
  { value: 'system:admin', label: '系统管理' }
]

const columns = [
  { colKey: 'id', title: 'ID', width: 80 },
  { colKey: 'name', title: '组名', width: 150 },
  { colKey: 'description', title: '描述', width: 200 },
  { colKey: 'permissions', title: '权限', width: 300 },
  { colKey: 'created_at', title: '创建时间', width: 160 },
  { colKey: 'operation', title: '操作', width: 150, fixed: 'right' }
]

// 获取权限标签
const getPermissionLabel = (permission: string) => {
  const found = availablePermissions.find(p => p.value === permission)
  return found ? found.label : permission
}

// 获取用户组列表
const fetchUserGroups = async () => {
  loading.value = true
  try {
    const response = await userGroupsApi.getUserGroups({
      page: pagination.current,
      limit: pagination.pageSize
    })
    userGroups.value = response.data || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取用户组列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 格式化日期
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString('zh-CN')
}

// 搜索
const handleSearch = () => {
  pagination.current = 1
  fetchUserGroups()
}

// 重置
const handleReset = () => {
  searchKeyword.value = ''
  pagination.current = 1
  fetchUserGroups()
}

// 分页变化
const handlePageChange = (pageInfo: any) => {
  pagination.current = pageInfo.current
  pagination.pageSize = pageInfo.pageSize
  fetchUserGroups()
}

// 编辑用户组
const handleEdit = (userGroup: UserGroup) => {
  editingUserGroup.value = userGroup
  Object.assign(formData, {
    name: userGroup.name,
    description: userGroup.description,
    permissions: [...userGroup.permissions]
  })
  showCreateDialog.value = true
}

// 删除用户组
const handleDelete = (userGroup: UserGroup) => {
  const dialog = DialogPlugin.confirm({
    header: '确认删除',
    body: `确定要删除用户组 "${userGroup.name}" 吗？此操作不可恢复。`,
    onConfirm: async () => {
      try {
        await userGroupsApi.deleteUserGroup(userGroup.id)
        MessagePlugin.success('删除成功')
        fetchUserGroups()
        dialog.destroy()
      } catch (error) {
        console.error('删除用户组失败:', error)
      }
    }
  })
}

// 提交表单
const handleSubmit = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  try {
    if (editingUserGroup.value) {
      // 编辑用户组
      await userGroupsApi.updateUserGroup(editingUserGroup.value.id, {
        name: formData.name,
        description: formData.description,
        permissions: formData.permissions
      })
      MessagePlugin.success('更新成功')
    } else {
      // 新增用户组
      await userGroupsApi.createUserGroup({
        name: formData.name,
        description: formData.description,
        permissions: formData.permissions
      })
      MessagePlugin.success('创建成功')
    }
    
    showCreateDialog.value = false
    fetchUserGroups()
  } catch (error) {
    console.error('操作失败:', error)
  }
}

// 取消操作
const handleCancel = () => {
  showCreateDialog.value = false
  editingUserGroup.value = null
  Object.assign(formData, {
    name: '',
    description: '',
    permissions: []
  })
}

onMounted(() => {
  fetchUserGroups()
})
</script>

<style scoped>
.user-groups-page {
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

.permissions-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 12px;
}

.permissions-grid .t-checkbox {
  margin: 0;
}
</style>
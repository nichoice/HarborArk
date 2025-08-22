import request from './request'
import type { UserGroup } from '@/types/auth'

export interface CreateUserGroupRequest {
  name: string
  description: string
  permissions: string[]
}

export interface UpdateUserGroupRequest {
  name?: string
  description?: string
  permissions?: string[]
}

export const userGroupsApi = {
  // 获取用户组列表
  getUserGroups: (params?: { page?: number; page_size?: number }) => {
    return request.get('/user-groups', { params })
  },

  // 获取单个用户组
  getUserGroup: (id: number) => {
    return request.get(`/user-groups/${id}`)
  },

  // 创建用户组
  createUserGroup: (data: CreateUserGroupRequest) => {
    return request.post('/user-groups', data)
  },

  // 更新用户组
  updateUserGroup: (id: number, data: UpdateUserGroupRequest) => {
    return request.put(`/user-groups/${id}`, data)
  },

  // 删除用户组
  deleteUserGroup: (id: number) => {
    return request.delete(`/user-groups/${id}`)
  }
}

import request from './request'
import type { User } from '@/types/auth'

export interface CreateUserRequest {
  username: string
  password: string
  email: string
  full_name: string
  user_group_id: number
}

export interface UpdateUserRequest {
  email?: string
  full_name?: string
  user_group_id?: number
  is_active?: boolean
}

export const usersApi = {
  // 获取用户列表
  getUsers: (params?: { page?: number; limit?: number }) => {
    return request.get('/users', { params })
  },

  // 获取单个用户
  getUser: (id: number) => {
    return request.get(`/users/${id}`)
  },

  // 创建用户
  createUser: (data: CreateUserRequest) => {
    return request.post('/users', data)
  },

  // 更新用户
  updateUser: (id: number, data: UpdateUserRequest) => {
    return request.put(`/users/${id}`, data)
  },

  // 删除用户
  deleteUser: (id: number) => {
    return request.delete(`/users/${id}`)
  }
}
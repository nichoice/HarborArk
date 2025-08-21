import request from './request'
import type { LoginRequest, LoginResponse } from '@/types/auth'

export const authApi = {
  // 用户登录
  login: (data: LoginRequest) => {
    return request.post<LoginResponse>('/auth/login', data)
  }
}
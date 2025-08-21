export interface LoginRequest {
  username: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
  message: string
}

export interface User {
  id: number
  username: string
  email: string
  full_name: string
  user_group_id: number
  user_group?: UserGroup
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface UserGroup {
  id: number
  name: string
  description: string
  permissions: string[]
  created_at: string
  updated_at: string
}
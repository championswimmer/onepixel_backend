import { post, get, patch } from './client'
import type { LoginRequest, UserResponse, UpdateUserRequest } from '../types'

export function loginUser(data: LoginRequest): Promise<UserResponse> {
  return post<UserResponse>('/users/login', data)
}

export function getUser(userId: number): Promise<UserResponse> {
  return get<UserResponse>(`/users/${userId}`)
}

export function updateUser(userId: number, data: UpdateUserRequest): Promise<UserResponse> {
  return patch<UserResponse>(`/users/${userId}`, data)
}

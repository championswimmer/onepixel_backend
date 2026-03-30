// Request types
export interface LoginRequest {
  email: string
  password: string
}

export interface CreateUrlRequest {
  long_url: string
}

export interface CreateUrlGroupRequest {
  short_path: string
  creator_id: number
}

export interface UpdateUserRequest {
  password: string
}

// Response types
export interface UserResponse {
  id: number
  email: string
  token: string
}

export interface UrlResponse {
  id: number
  short_url: string
  long_url: string
  creator_id: number
}

export interface UrlInfoResponse {
  long_url: string
  hit_count: number
}

export interface UrlGroupResponse {
  short_path: string
  creator_id: number
}

export interface ErrorResponse {
  status: number
  message: string
}

export interface StatsResponse {
  urls_count: number
  total_hits: number
}

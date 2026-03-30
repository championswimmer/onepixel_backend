import { get, post, put } from './client'
import type { CreateUrlRequest, CreateUrlGroupRequest, UrlResponse, UrlInfoResponse, UrlGroupResponse } from '../types'

export function getUrls(): Promise<UrlResponse[]> {
  return get<UrlResponse[]>('/urls')
}

export function createRandomUrl(data: CreateUrlRequest): Promise<UrlResponse> {
  return post<UrlResponse>('/urls', data)
}

export function createCustomUrl(shortcode: string, data: CreateUrlRequest): Promise<UrlResponse> {
  return put<UrlResponse>(`/urls/${shortcode}`, data)
}

export function createGroupedRandomUrl(group: string, data: CreateUrlRequest): Promise<UrlResponse> {
  return post<UrlResponse>(`/urls/groups/${group}/shorten`, data)
}

export function createGroupedCustomUrl(group: string, shortcode: string, data: CreateUrlRequest): Promise<UrlResponse> {
  return post<UrlResponse>(`/urls/groups/${group}/shorten/${shortcode}`, data)
}

export function createUrlGroup(data: CreateUrlGroupRequest): Promise<UrlGroupResponse> {
  return post<UrlGroupResponse>('/urls/groups', data)
}

export function getUrlGroups(): Promise<UrlGroupResponse[]> {
  return get<UrlGroupResponse[]>('/urls/groups')
}

export function getUrlInfo(shortcode: string): Promise<UrlInfoResponse> {
  return get<UrlInfoResponse>(`/urls/${shortcode}`)
}

export function getGroupedUrlInfo(group: string, shortcode: string): Promise<UrlInfoResponse> {
  return get<UrlInfoResponse>(`/urls/groups/${group}/${shortcode}`)
}

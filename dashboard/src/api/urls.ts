import { get, post, put } from './client'
import type { CreateUrlRequest, UrlResponse, UrlInfoResponse } from '../types'

export function getUrls(): Promise<UrlResponse[]> {
  return get<UrlResponse[]>('/urls')
}

export function createRandomUrl(data: CreateUrlRequest): Promise<UrlResponse> {
  return post<UrlResponse>('/urls', data)
}

export function createCustomUrl(shortcode: string, data: CreateUrlRequest): Promise<UrlResponse> {
  return put<UrlResponse>(`/urls/${shortcode}`, data)
}

export function getUrlInfo(shortcode: string): Promise<UrlInfoResponse> {
  return get<UrlInfoResponse>(`/urls/${shortcode}`)
}

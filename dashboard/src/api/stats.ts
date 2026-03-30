import { get } from './client'

export function getStats(): Promise<unknown> {
  return get<unknown>('/stats')
}

export function getUrlStats(shortcode: string): Promise<unknown> {
  return get<unknown>(`/stats/${shortcode}`)
}

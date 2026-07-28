import type { CompareResult, HistoryResult, Report, RivalSuggestion } from './types'

const API_BASE = '/api'

async function fetchJSON<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`)
  if (!res.ok) {
    const body = await res.json().catch(() => ({}))
    throw new Error((body as { error?: string }).error ?? `API error ${res.status}`)
  }
  return res.json() as Promise<T>
}

export function analyzeRepo(repo: string): Promise<Report> {
  return fetchJSON(`/analyze?repo=${encodeURIComponent(repo)}`)
}

export function getHistory(repo: string, period = 'weekly'): Promise<HistoryResult> {
  return fetchJSON(`/history?repo=${encodeURIComponent(repo)}&period=${period}`)
}

export function compareRepos(repo1: string, repo2: string): Promise<CompareResult> {
  return fetchJSON(
    `/compare?repo1=${encodeURIComponent(repo1)}&repo2=${encodeURIComponent(repo2)}`,
  )
}

export function getRivals(repo: string): Promise<RivalSuggestion[]> {
  return fetchJSON(`/rivals?repo=${encodeURIComponent(repo)}`)
}

export function healthCheck(): Promise<{ status: string; app: string }> {
  return fetchJSON('/health')
}

export type Status = 'OK' | 'WARN' | 'FAIL'

export interface CheckResult {
  name: string
  score: number
  maxScore: number
  status: Status
  detail: string
  suggestion: string
}

export interface Category {
  name: string
  score: number
  maxScore: number
  checks: CheckResult[]
}

export interface Report {
  owner: string
  repo: string
  fullName: string
  score: number
  maxScore: number
  categories: Category[]
  analyzedAt: string
}

export interface HistoryPoint {
  date: string
  score: number
  event?: string
}

export interface HistoryResult {
  points: HistoryPoint[]
  coveredDays: number
  truncated: boolean
}

export interface CompareResult {
  repo1: Report
  repo2: Report
  diff: Record<string, {
    repo1Score: number
    repo2Score: number
    delta: number
    winner: string
  }>
}

export interface RivalSuggestion {
  fullName: string
  reason: string
  similarity: string
}

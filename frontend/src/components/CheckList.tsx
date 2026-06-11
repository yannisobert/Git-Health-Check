import type { CheckResult } from '../api/types'

const STATUS_CONFIG = {
  ok:   { char: '✓', bg: '#ecfdf5', color: '#059669' },
  warn: { char: '~', bg: '#fffbeb', color: '#d97706' },
  fail: { char: '✕', bg: '#fef2f2', color: '#dc2626' },
} as const

interface CheckListProps {
  checks: CheckResult[]
  compact?: boolean
}

export function CheckList({ checks, compact }: CheckListProps) {
  return (
    <div className="check-list">
      {checks.map((check) => {
        const status = check.status.toLowerCase() as keyof typeof STATUS_CONFIG
        const cfg = STATUS_CONFIG[status] ?? STATUS_CONFIG.fail
        return (
          <div key={check.name} className={`check-item check-${status}${compact ? ' check-compact' : ''}`}>
            <div className="check-icon" style={{ background: cfg.bg, color: cfg.color }}>
              {cfg.char}
            </div>
            <div className="check-body">
              <div className="check-header">
                <strong className="check-name">{check.name}</strong>
                <span className="check-score">{check.score}/{check.maxScore}</span>
              </div>
              {check.detail && <p className="check-detail">{check.detail}</p>}
              {check.suggestion && (
                <p className="suggestion">
                  <span>→</span> {check.suggestion}
                </p>
              )}
            </div>
          </div>
        )
      })}
    </div>
  )
}

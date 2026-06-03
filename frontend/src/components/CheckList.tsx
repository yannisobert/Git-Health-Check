import type { CheckResult } from '../api/types'

interface CheckListProps {
  checks: CheckResult[]
}

export function CheckList({ checks }: CheckListProps) {
  return (
    <ul className="check-list">
      {checks.map((check) => (
        <li key={check.name} className={`check-item check-${check.status.toLowerCase()}`}>
          <div className="check-header">
            <strong>{check.name}</strong>
            <span>
              {check.score}/{check.maxScore}
            </span>
          </div>
          {check.detail && <p>{check.detail}</p>}
          {check.suggestion && <p className="suggestion">{check.suggestion}</p>}
        </li>
      ))}
    </ul>
  )
}

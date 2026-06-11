interface StickyScoreBarProps {
  repo: string
  score: number
  maxScore: number
  visible: boolean
}

function scoreColor(pct: number): string {
  if (pct >= 70) return '#059669'
  if (pct >= 40) return '#d97706'
  return '#dc2626'
}

function scoreLabel(pct: number): string {
  if (pct >= 70) return 'Healthy'
  if (pct >= 40) return 'Needs attention'
  return 'Critical'
}

export function StickyScoreBar({ repo, score, maxScore, visible }: StickyScoreBarProps) {
  const pct = maxScore > 0 ? Math.round((score / maxScore) * 100) : 0
  const color = scoreColor(pct)
  const r = 11
  const circumference = 2 * Math.PI * r
  const dash = (pct / 100) * circumference

  return (
    <div className={`sticky-score-bar${visible ? ' visible' : ''}`}>
      <div className="sticky-score-inner">
        <svg viewBox="0 0 28 28" width="28" height="28" style={{ flexShrink: 0 }}>
          <circle cx="14" cy="14" r={r} fill="none" stroke={`${color}30`} strokeWidth="3" />
          <circle
            cx="14" cy="14" r={r}
            fill="none"
            stroke={color}
            strokeWidth="3"
            strokeLinecap="round"
            strokeDasharray={circumference}
            strokeDashoffset={circumference - dash}
            transform="rotate(-90 14 14)"
            style={{ filter: `drop-shadow(0 0 3px ${color}99)` }}
          />
        </svg>
        <span className="sticky-score-pct" style={{ color }}>{pct}</span>
        <span className="sticky-score-sep" />
        <span className="sticky-score-repo">{repo}</span>
        <span className="sticky-score-label" style={{ color, background: `${color}18`, border: `1px solid ${color}35` }}>
          {scoreLabel(pct)}
        </span>
      </div>
    </div>
  )
}

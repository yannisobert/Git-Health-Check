import type { Category } from '../api/types'

interface ScoreCardProps {
  score: number
  maxScore: number
  repo: string
  categories?: Category[]
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

export function ScoreCard({ score, maxScore, repo, categories }: ScoreCardProps) {
  const pct = maxScore > 0 ? Math.round((score / maxScore) * 100) : 0
  const color = scoreColor(pct)
  const r = 42
  const circumference = 2 * Math.PI * r
  const dash = (pct / 100) * circumference

  return (
    <div
      className="score-card"
      style={{ boxShadow: `0 0 0 1px ${color}30, 0 8px 48px ${color}18` }}
    >
      <div className="score-card-hero">
        <div className="score-card-top">
          <div className="score-ring-wrap">
            <svg viewBox="0 0 100 100" width="130" height="130">
              <circle cx="50" cy="50" r={r} fill="none" stroke="rgba(255,255,255,0.14)" strokeWidth="7" />
              <circle
                cx="50" cy="50" r={r}
                fill="none"
                stroke={color}
                strokeWidth="7"
                strokeLinecap="round"
                strokeDasharray={circumference}
                strokeDashoffset={circumference - dash}
                transform="rotate(-90 50 50)"
                style={{
                  transition: 'stroke-dashoffset 0.7s ease',
                  filter: `drop-shadow(0 0 7px ${color})`,
                }}
              />
              <text
                x="50" y="50"
                textAnchor="middle"
                dominantBaseline="middle"
                fontSize="24"
                fontWeight="700"
                fill="white"
                fontFamily="Inter, system-ui, sans-serif"
              >
                {pct}
              </text>
            </svg>
          </div>

          <div className="score-card-info">
            <h2 className="score-repo">{repo}</h2>
            <p className="score-label">{scoreLabel(pct)}</p>
          </div>
        </div>
      </div>

      {categories && (
        <div className="category-bars">
          {categories.map((cat) => {
            const catPct = cat.maxScore > 0 ? (cat.score / cat.maxScore) * 100 : 0
            return (
              <div key={cat.name} className="category-bar-row">
                <span className="category-name">{cat.name}</span>
                <div className="bar-track">
                  <div
                    className="bar-fill"
                    style={{ width: `${catPct}%`, background: scoreColor(catPct) }}
                  />
                </div>
                <span className="category-score">{cat.score}/{cat.maxScore}</span>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}

import type { Category } from '../api/types'

interface ScoreCardProps {
  score: number
  maxScore: number
  repo: string
  categories?: Category[]
}

function scoreColor(pct: number): string {
  if (pct >= 70) return '#16a34a'
  if (pct >= 40) return '#d97706'
  return '#dc2626'
}

export function ScoreCard({ score, maxScore, repo, categories }: ScoreCardProps) {
  const pct = maxScore > 0 ? Math.round((score / maxScore) * 100) : 0
  const color = scoreColor(pct)

  return (
    <div className="score-card">
      <div className="score-card-header">
        <h2>{repo}</h2>
        <div className="score-circle" style={{ color }}>
          <span className="score-number">{pct}</span>
          <span className="score-slash">/100</span>
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
                <span className="category-score">
                  {cat.score}/{cat.maxScore}
                </span>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}

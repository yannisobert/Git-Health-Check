import { useEffect, useState } from 'react'
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
  const targetDash = (pct / 100) * circumference

  const [dashOffset, setDashOffset] = useState(circumference)
  useEffect(() => {
    const t = setTimeout(() => setDashOffset(circumference - targetDash), 60)
    return () => clearTimeout(t)
  }, [circumference, targetDash])

  const [displayPct, setDisplayPct] = useState(0)
  useEffect(() => {
    const duration = 900
    const start = performance.now()
    let raf: number
    const step = (now: number) => {
      const t = Math.min((now - start) / duration, 1)
      const eased = 1 - Math.pow(1 - t, 3)
      setDisplayPct(Math.round(eased * pct))
      if (t < 1) raf = requestAnimationFrame(step)
    }
    raf = requestAnimationFrame(step)
    return () => cancelAnimationFrame(raf)
  }, [pct])

  const [barWidths, setBarWidths] = useState<number[]>([])
  useEffect(() => {
    if (!categories) return
    const t = setTimeout(() => {
      setBarWidths(categories.map(cat => cat.maxScore > 0 ? (cat.score / cat.maxScore) * 100 : 0))
    }, 250)
    return () => clearTimeout(t)
  }, [categories])

  return (
    <div
      className="score-card"
      style={{
        background: `linear-gradient(135deg, #ffffff 40%, ${color}18 100%)`,
        border: `1px solid ${color}40`,
        boxShadow: `0 4px 32px ${color}18`,
      }}
    >
      <div className="score-card-top">
        <div className="score-ring-wrap">
          <svg viewBox="0 0 100 100" width="130" height="130">
            <circle cx="50" cy="50" r={r} fill="none" stroke={`${color}25`} strokeWidth="7" />
            <circle
              cx="50" cy="50" r={r}
              fill="none"
              stroke={color}
              strokeWidth="7"
              strokeLinecap="round"
              strokeDasharray={circumference}
              strokeDashoffset={dashOffset}
              transform="rotate(-90 50 50)"
              style={{
                transition: 'stroke-dashoffset 0.9s cubic-bezier(0.34, 1.56, 0.64, 1)',
                filter: `drop-shadow(0 0 6px ${color}99)`,
              }}
            />
            <text
              x="50" y="50"
              textAnchor="middle"
              dominantBaseline="middle"
              fontSize="24"
              fontWeight="700"
              fill={color}
              fontFamily="Plus Jakarta Sans, system-ui, sans-serif"
            >
              {displayPct}
            </text>
          </svg>
        </div>

        <div className="score-card-info">
          <h2 className="score-repo">{repo}</h2>
          <p className="score-label" style={{ color, background: `${color}18`, border: `1px solid ${color}35` }}>
            {scoreLabel(pct)}
          </p>
        </div>
      </div>

      {categories && (
        <div className="category-bars">
          {categories.map((cat, i) => {
            const catPct = cat.maxScore > 0 ? (cat.score / cat.maxScore) * 100 : 0
            const animWidth = barWidths[i] ?? 0
            return (
              <div key={cat.name} className="category-bar-row">
                <span className="category-name">{cat.name}</span>
                <div className="bar-track">
                  <div
                    className="bar-fill"
                    style={{
                      width: `${animWidth}%`,
                      background: scoreColor(catPct),
                      transition: `width 0.6s cubic-bezier(0.22, 1, 0.36, 1) ${i * 80}ms`,
                    }}
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

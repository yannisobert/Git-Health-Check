import type { HistoryPoint } from '../api/types'

interface EvolutionChartProps {
  points: HistoryPoint[]
}

export function EvolutionChart({ points }: EvolutionChartProps) {
  if (points.length === 0) {
    return <p className="placeholder">Evolution chart — available after feature/history</p>
  }

  return (
    <div className="evolution-chart">
      <p>Score history ({points.length} points)</p>
      <ul>
        {points.map((p) => (
          <li key={p.date}>
            {new Date(p.date).toLocaleDateString()} — {p.score}
            {p.event ? ` (${p.event})` : ''}
          </li>
        ))}
      </ul>
    </div>
  )
}

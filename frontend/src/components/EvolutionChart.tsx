import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  LinearScale,
  LineElement,
  PointElement,
  Tooltip,
} from 'chart.js'
import { useEffect, useRef, useState } from 'react'
import { Line } from 'react-chartjs-2'
import type { HistoryPoint } from '../api/types'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Filler)

type Period = 'weekly' | 'monthly'

interface EvolutionChartProps {
  points: HistoryPoint[]
  coveredDays: number
  truncated: boolean
  period: Period
  onPeriodChange: (p: Period) => void
  loading?: boolean
}

const PERIODS: { value: Period; label: string }[] = [
  { value: 'weekly', label: '6M' },
  { value: 'monthly', label: '1Y' },
]

function formatCoverage(days: number) {
  return days < 30 ? `${days}d` : `${Math.round(days / 30)}mo`
}

function scoreColor(pct: number) {
  if (pct >= 70) return '#059669'
  if (pct >= 40) return '#d97706'
  return '#dc2626'
}

function hexToRgb(hex: string) {
  const r = parseInt(hex.slice(1, 3), 16)
  const g = parseInt(hex.slice(3, 5), 16)
  const b = parseInt(hex.slice(5, 7), 16)
  return `${r},${g},${b}`
}

const CHART_ANIMATION_MS = 900

export function EvolutionChart({ points, coveredDays, truncated, period, onPeriodChange, loading }: EvolutionChartProps) {
  const [showSpinner, setShowSpinner] = useState(false)
  const timerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  useEffect(() => {
    if (loading) {
      if (timerRef.current) clearTimeout(timerRef.current)
      setShowSpinner(true)
    } else {
      timerRef.current = setTimeout(() => setShowSpinner(false), CHART_ANIMATION_MS)
    }
    return () => { if (timerRef.current) clearTimeout(timerRef.current) }
  }, [loading])
  const labels = points.map((p) => {
    const d = new Date(p.date)
    if (period === 'monthly') {
      return d.toLocaleDateString('en-US', { month: 'short' })
    }
    return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  })

  const scores = points.map((p) => p.score)
  const maxScore = scores.length ? Math.max(...scores) : null
  const minScore = scores.length ? Math.min(...scores) : null
  const maxIdx = scores.indexOf(maxScore ?? 0)
  const minIdx = scores.indexOf(minScore ?? 0)

  // Dynamic color based on last score
  const lastScore = scores[scores.length - 1] ?? 0
  const lineColor = scoreColor(lastScore)
  const lineRgb = hexToRgb(lineColor)

  const data = {
    labels,
    datasets: [
      {
        label: 'Score',
        data: scores,
        borderColor: lineColor,
        backgroundColor: `rgba(${lineRgb},0.07)`,
        fill: true,
        tension: 0.35,
        pointRadius: points.map((p, i) =>
          p.event || i === maxIdx || i === minIdx ? 5 : 2,
        ),
        pointBackgroundColor: points.map((p, i) => {
          if (p.event) return '#d97706'
          if (i === maxIdx) return lineColor
          if (i === minIdx) return '#dc2626'
          return lineColor
        }),
        pointBorderColor: points.map((p, i) => {
          if (p.event) return '#d97706'
          if (i === maxIdx || i === minIdx) return '#ffffff'
          return lineColor
        }),
        pointBorderWidth: points.map((_p, i) =>
          i === maxIdx || i === minIdx ? 2 : 0,
        ),
      },
    ],
  }

  const options = {
    responsive: true,
    animation: { duration: 900, easing: 'easeInOutQuart' as const },
    scales: {
      y: {
        min: 0,
        max: 100,
        grid: { color: '#f3f4f6' },
        ticks: { color: '#94a3b8', font: { size: 11 } },
      },
      x: {
        grid: { display: false },
        ticks: { color: '#94a3b8', font: { size: 11 } },
      },
    },
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          label: (ctx: { dataIndex: number; parsed: { y: number | null } }) => {
            const pct = ctx.parsed.y ?? 0
            const prev = scores[ctx.dataIndex - 1]
            const delta = prev != null ? pct - prev : null
            const deltaStr =
              delta != null
                ? delta > 0
                  ? ` (+${delta} pts)`
                  : delta < 0
                    ? ` (${delta} pts)`
                    : ' (=)'
                : ''
            return `Score: ${pct}${deltaStr}`
          },
          afterLabel: (ctx: { dataIndex: number }) => {
            const event = points[ctx.dataIndex]?.event
            return event ? `Event: ${event}` : ''
          },
        },
      },
    },
  }

  const minDate = minIdx >= 0 ? labels[minIdx] : null
  const maxDate = maxIdx >= 0 ? labels[maxIdx] : null

  return (
    <div className="evolution-chart">
      <div className="evolution-header">
        <div className="evolution-title-group">
          <h3>Score History</h3>
          {scores.length > 1 && minScore !== null && maxScore !== null && (
            <div className="evolution-stats">
              <span className="evolution-stat evolution-stat--max" title={maxDate ?? ''}>
                ↑ {maxScore}
              </span>
              <span className="evolution-stat evolution-stat--min" title={minDate ?? ''}>
                ↓ {minScore}
              </span>
            </div>
          )}
        </div>
        <div className="period-tabs">
          {PERIODS.map((p) => (
            <button
              key={p.value}
              className={`period-tab${period === p.value ? ' active' : ''}`}
              onClick={() => onPeriodChange(p.value)}
            >
              {p.label}
            </button>
          ))}
        </div>
      </div>

      {!loading && points.length > 0 && (
        <p className={`coverage-note${truncated ? ' truncated' : ''}`}>
          Showing {formatCoverage(coveredDays)} of history
          {truncated && ' — GitHub fetch limit reached, actual history may be longer'}
        </p>
      )}

      <div className="chart-wrapper">
        {showSpinner && <div className="chart-loading"><span className="chart-spinner" /></div>}
        {!loading && points.length === 0 ? (
          <p className="chart-empty">Not enough history yet.</p>
        ) : (
          <Line data={data} options={options} />
        )}
      </div>
    </div>
  )
}

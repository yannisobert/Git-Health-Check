import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  LinearScale,
  LineElement,
  PointElement,
  Tooltip,
} from 'chart.js'
import { Line } from 'react-chartjs-2'
import type { HistoryPoint } from '../api/types'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Filler)

type Period = 'weekly' | 'monthly'

interface EvolutionChartProps {
  points: HistoryPoint[]
  period: Period
  onPeriodChange: (p: Period) => void
}

const PERIODS: { value: Period; label: string }[] = [
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: '1Y' },
]

export function EvolutionChart({ points, period, onPeriodChange }: EvolutionChartProps) {
  const labels = points.map((p) =>
    new Date(p.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
  )

  const data = {
    labels,
    datasets: [
      {
        label: 'Score',
        data: points.map((p) => p.score),
        borderColor: '#6366f1',
        backgroundColor: 'rgba(99,102,241,0.07)',
        fill: true,
        tension: 0.35,
        pointRadius: points.map((p) => (p.event ? 5 : 2)),
        pointBackgroundColor: points.map((p) => (p.event ? '#d97706' : '#6366f1')),
      },
    ],
  }

  const options = {
    responsive: true,
    scales: {
      y: { min: 0, max: 100, grid: { color: '#f3f4f6' }, ticks: { color: '#94a3b8', font: { size: 11 } } },
      x: { grid: { display: false }, ticks: { color: '#94a3b8', font: { size: 11 } } },
    },
    plugins: {
      legend: { display: false },
      tooltip: {
        callbacks: {
          afterLabel: (_ctx: { dataIndex: number }) => {
            const event = points[_ctx.dataIndex]?.event
            return event ? `Event: ${event}` : ''
          },
        },
      },
    },
  }

  return (
    <div className="evolution-chart">
      <div className="evolution-header">
        <h3>Score History</h3>
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

      {points.length === 0 ? (
        <p className="chart-empty">Not enough history yet.</p>
      ) : (
        <Line data={data} options={options} />
      )}
    </div>
  )
}

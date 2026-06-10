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

interface EvolutionChartProps {
  points: HistoryPoint[]
}

export function EvolutionChart({ points }: EvolutionChartProps) {
  if (points.length === 0) return null

  const labels = points.map((p) =>
    new Date(p.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }),
  )

  const data = {
    labels,
    datasets: [
      {
        label: 'Score',
        data: points.map((p) => p.score),
        borderColor: '#2563eb',
        backgroundColor: 'rgba(37,99,235,0.08)',
        fill: true,
        tension: 0.3,
        pointRadius: points.map((p) => (p.event ? 5 : 2)),
        pointBackgroundColor: points.map((p) => (p.event ? '#d97706' : '#2563eb')),
      },
    ],
  }

  const options = {
    responsive: true,
    scales: {
      y: { min: 0, max: 100, grid: { color: '#f3f4f6' } },
      x: { grid: { display: false } },
    },
    plugins: {
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
      <h3>Score History</h3>
      <Line data={data} options={options} />
    </div>
  )
}

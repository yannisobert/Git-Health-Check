import {
  CategoryScale,
  Chart as ChartJS,
  Filler,
  LinearScale,
  LineElement,
  PointElement,
  Tooltip,
  Legend,
} from 'chart.js'
import { Line } from 'react-chartjs-2'
import type { HistoryPoint } from '../api/types'

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Tooltip, Filler, Legend)

interface CompareChartProps {
  points1: HistoryPoint[]
  points2: HistoryPoint[]
  label1: string
  label2: string
}

export function CompareChart({ points1, points2, label1, label2 }: CompareChartProps) {
  const labels = points1.length >= points2.length
    ? points1.map((p) => new Date(p.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }))
    : points2.map((p) => new Date(p.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' }))

  const data = {
    labels,
    datasets: [
      {
        label: label1,
        data: points1.map((p) => p.score),
        borderColor: '#2563eb',
        backgroundColor: 'rgba(37,99,235,0.06)',
        fill: true,
        tension: 0.3,
        pointRadius: 2,
      },
      {
        label: label2,
        data: points2.map((p) => p.score),
        borderColor: '#d97706',
        backgroundColor: 'rgba(217,119,6,0.06)',
        fill: true,
        tension: 0.3,
        pointRadius: 2,
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
      legend: { position: 'top' as const },
    },
  }

  return (
    <div className="evolution-chart">
      <h3>Score Evolution</h3>
      <Line data={data} options={options} />
    </div>
  )
}

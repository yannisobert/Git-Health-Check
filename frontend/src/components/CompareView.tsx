import type { CompareResult, HistoryPoint } from '../api/types'
import { CompareChart } from './CompareChart'
import { ScoreCard } from './ScoreCard'

interface CompareViewProps {
  result: CompareResult
  history1: HistoryPoint[]
  history2: HistoryPoint[]
}

export function CompareView({ result, history1, history2 }: CompareViewProps) {
  const { repo1, repo2, diff } = result

  return (
    <div className="compare-view">
      <div className="compare-grid">
        <ScoreCard
          repo={repo1.fullName}
          score={repo1.score}
          maxScore={repo1.maxScore}
          categories={repo1.categories}
        />
        <ScoreCard
          repo={repo2.fullName}
          score={repo2.score}
          maxScore={repo2.maxScore}
          categories={repo2.categories}
        />
      </div>

      {(history1.length > 0 || history2.length > 0) && (
        <CompareChart
          points1={history1}
          points2={history2}
          label1={repo1.fullName}
          label2={repo2.fullName}
        />
      )}

      <div className="diff-table">
        <h3>Check by check</h3>
        <table>
          <thead>
            <tr>
              <th>Check</th>
              <th>{repo1.repo}</th>
              <th>{repo2.repo}</th>
              <th>Winner</th>
            </tr>
          </thead>
          <tbody>
            {Object.entries(diff).map(([name, d]) => (
              <tr key={name}>
                <td>{name}</td>
                <td className={d.winner === repo1.fullName ? 'cell-winner' : ''}>{d.repo1Score}</td>
                <td className={d.winner === repo2.fullName ? 'cell-winner' : ''}>{d.repo2Score}</td>
                <td className="cell-winner-name">
                  {d.winner === 'tie' ? '—' : d.winner.split('/')[1]}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

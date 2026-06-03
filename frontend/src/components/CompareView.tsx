import type { CompareResult } from '../api/types'
import { ScoreCard } from './ScoreCard'

interface CompareViewProps {
  result: CompareResult
}

export function CompareView({ result }: CompareViewProps) {
  return (
    <div className="compare-view">
      <div className="compare-grid">
        <ScoreCard
          repo={result.repo1.fullName}
          score={result.repo1.score}
          maxScore={result.repo1.maxScore}
        />
        <ScoreCard
          repo={result.repo2.fullName}
          score={result.repo2.score}
          maxScore={result.repo2.maxScore}
        />
      </div>
    </div>
  )
}

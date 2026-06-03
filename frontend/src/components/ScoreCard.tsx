interface ScoreCardProps {
  score: number
  maxScore: number
  repo: string
}

export function ScoreCard({ score, maxScore, repo }: ScoreCardProps) {
  const pct = maxScore > 0 ? Math.round((score / maxScore) * 100) : 0

  return (
    <div className="score-card">
      <h2>{repo}</h2>
      <div className="score-value">{pct}</div>
      <p className="score-detail">
        {score} / {maxScore} points
      </p>
    </div>
  )
}

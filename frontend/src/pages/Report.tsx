import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { analyzeRepo, getHistory, getRivals } from '../api/client'
import type { HistoryPoint, Report as ReportType, RivalSuggestion } from '../api/types'
import { CheckList } from '../components/CheckList'
import { EvolutionChart } from '../components/EvolutionChart'
import { RivalSuggestions } from '../components/RivalSuggestions'
import { ScoreCard } from '../components/ScoreCard'

export function Report() {
  const { owner = '', repo = '' } = useParams()
  const fullName = `${owner}/${repo}`

  const [report, setReport] = useState<ReportType | null>(null)
  const [history, setHistory] = useState<HistoryPoint[]>([])
  const [rivals, setRivals] = useState<RivalSuggestion[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setLoading(true)
    setError(null)

    analyzeRepo(fullName)
      .then((r) => {
        setReport(r)
        return Promise.all([
          getHistory(fullName).then(setHistory).catch(() => {}),
          getRivals(fullName).then(setRivals).catch(() => {}),
        ])
      })
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [fullName])

  if (loading) {
    return (
      <section className="page report">
        <h1>Analyzing {fullName}…</h1>
        <div className="spinner" />
      </section>
    )
  }

  if (error) {
    return (
      <section className="page report">
        <h1>Report — {fullName}</h1>
        <p className="error-msg">{error}</p>
      </section>
    )
  }

  if (!report) return null

  const allChecks = report.categories.flatMap((c) => c.checks)

  return (
    <section className="page report">
      <ScoreCard
        repo={report.fullName}
        score={report.score}
        maxScore={report.maxScore}
        categories={report.categories}
      />
      <CheckList checks={allChecks} />
      <EvolutionChart points={history} />
      <RivalSuggestions rivals={rivals} />
    </section>
  )
}

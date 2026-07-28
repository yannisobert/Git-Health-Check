import { useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import { analyzeRepo, getHistory, getRivals } from '../api/client'
import type { HistoryResult, Report as ReportType, RivalSuggestion } from '../api/types'
import { AnalyzingLoader } from '../components/AnalyzingLoader'
import { CheckList } from '../components/CheckList'
import { EvolutionChart } from '../components/EvolutionChart'
import { RivalSuggestions } from '../components/RivalSuggestions'
import { ScoreCard } from '../components/ScoreCard'
import { StickyScoreBar } from '../components/StickyScoreBar'

type Period = 'weekly' | 'monthly'

const CATEGORY_ICONS: Record<string, string> = {
  Maintenance: '🔧',
  Activity: '⚡',
  Conventions: '📐',
  'CI/CD': '🚀',
}

function scoreColor(pct: number) {
  if (pct >= 70) return '#059669'
  if (pct >= 40) return '#d97706'
  return '#dc2626'
}

export function Report() {
  const { owner = '', repo = '' } = useParams()
  const fullName = `${owner}/${repo}`

  const [report, setReport] = useState<ReportType | null>(null)
  const [history, setHistory] = useState<HistoryResult>({ points: [], coveredDays: 0, truncated: false })
  const [rivals, setRivals] = useState<RivalSuggestion[]>([])
  const [period, setPeriod] = useState<Period>('weekly')
  const [loading, setLoading] = useState(true)
  const [historyLoading, setHistoryLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [stickyVisible, setStickyVisible] = useState(false)
  const cardRefs = useRef<(HTMLDivElement | null)[]>([])
  const scoreCardRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    setLoading(true)
    setError(null)
    Promise.all([
      analyzeRepo(fullName),
      getRivals(fullName).catch((): RivalSuggestion[] => []),
    ])
      .then(([r, riv]) => {
        setReport(r)
        setRivals(riv)
      })
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [fullName])

  useEffect(() => {
    setHistoryLoading(true)
    getHistory(fullName, period)
      .then(setHistory)
      .catch(() => setHistory({ points: [], coveredDays: 0, truncated: false }))
      .finally(() => setHistoryLoading(false))
  }, [fullName, period])

  // Sticky bar: show when ScoreCard leaves viewport
  useEffect(() => {
    if (!scoreCardRef.current) return
    const obs = new IntersectionObserver(
      ([entry]) => setStickyVisible(!entry.isIntersecting),
      { threshold: 0, rootMargin: '-60px 0px 0px 0px' }
    )
    obs.observe(scoreCardRef.current)
    return () => obs.disconnect()
  }, [report])

  // Staggered scroll reveal for category cards
  useEffect(() => {
    if (!report) return
    const observers: IntersectionObserver[] = []
    cardRefs.current.forEach((el, i) => {
      if (!el) return
      const obs = new IntersectionObserver(
        ([entry]) => {
          if (entry.isIntersecting) {
            setTimeout(() => el.classList.add('visible'), i * 70)
            obs.disconnect()
          }
        },
        { threshold: 0.1 }
      )
      obs.observe(el)
      observers.push(obs)
    })
    return () => observers.forEach(o => o.disconnect())
  }, [report])

  if (loading) return <AnalyzingLoader repo={fullName} />

  const banner = (
    <div className="page-banner">
      <div className="page-banner-inner">
        <p className="page-banner-eyebrow">Health Report</p>
        <h1>{fullName}</h1>
      </div>
    </div>
  )

  if (error) {
    return (
      <>
        {banner}
        <div className="container-overlap">
          <p className="error-msg">{error}</p>
        </div>
      </>
    )
  }

  if (!report) return null

  return (
    <>
      {banner}
      {report && (
        <StickyScoreBar
          repo={report.fullName}
          score={report.score}
          maxScore={report.maxScore}
          visible={stickyVisible}
        />
      )}
      <div className="container-overlap">
        <div ref={scoreCardRef}>
          <ScoreCard
            repo={report.fullName}
            score={report.score}
            maxScore={report.maxScore}
            categories={report.categories}
          />
        </div>

        <div className="checks-grid">
          {report.categories.map((cat, i) => {
            const pct = cat.maxScore > 0 ? Math.round((cat.score / cat.maxScore) * 100) : 0
            const color = scoreColor(pct)
            const icon = CATEGORY_ICONS[cat.name] ?? '📋'
            return (
              <div
                key={cat.name}
                className="cat-card"
                ref={el => { cardRefs.current[i] = el }}
                style={{ borderTop: `2px solid ${color}` }}
              >
                <div className="cat-card-header">
                  <span className="cat-card-title">
                    <span className="cat-card-icon">{icon}</span>
                    {cat.name}
                  </span>
                  <span className="cat-card-score" style={{ color, background: `${color}18` }}>
                    {cat.score}/{cat.maxScore}
                  </span>
                </div>
                <CheckList checks={cat.checks} compact />
              </div>
            )
          })}
        </div>

        <EvolutionChart
          points={history.points}
          coveredDays={history.coveredDays}
          truncated={history.truncated}
          period={period}
          onPeriodChange={setPeriod}
          loading={historyLoading}
        />
        <RivalSuggestions rivals={rivals} />
      </div>
    </>
  )
}

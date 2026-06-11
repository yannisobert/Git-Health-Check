import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { compareRepos, getHistory } from '../api/client'
import type { CompareResult, HistoryPoint } from '../api/types'
import { AnalyzingLoader } from '../components/AnalyzingLoader'
import { CompareView } from '../components/CompareView'

export function Compare() {
  const { o1 = '', r1 = '', o2 = '', r2 = '' } = useParams()
  const repo1 = `${o1}/${r1}`
  const repo2 = `${o2}/${r2}`

  const [result, setResult] = useState<CompareResult | null>(null)
  const [history1, setHistory1] = useState<HistoryPoint[]>([])
  const [history2, setHistory2] = useState<HistoryPoint[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    setLoading(true)
    setError(null)

    Promise.all([
      compareRepos(repo1, repo2),
      getHistory(repo1).catch((): HistoryPoint[] => []),
      getHistory(repo2).catch((): HistoryPoint[] => []),
    ])
      .then(([compareResult, h1, h2]) => {
        setResult(compareResult)
        setHistory1(h1)
        setHistory2(h2)
      })
      .catch((e: Error) => setError(e.message))
      .finally(() => setLoading(false))
  }, [repo1, repo2])

  if (loading) return <AnalyzingLoader repo={`${repo1} · ${repo2}`} />

  const banner = (
    <div className="page-banner">
      <div className="page-banner-inner">
        <p className="page-banner-eyebrow">Comparison</p>
        <div className="banner-compare">
          <span>{repo1}</span>
          <span className="vs-badge">VS</span>
          <span>{repo2}</span>
        </div>
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

  if (!result) return null

  return (
    <>
      {banner}
      <div className="container-overlap">
        <CompareView result={result} history1={history1} history2={history2} />
      </div>
    </>
  )
}

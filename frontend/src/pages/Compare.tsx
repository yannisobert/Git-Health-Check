import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { compareRepos, getHistory } from '../api/client'
import type { CompareResult, HistoryPoint } from '../api/types'
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

  if (loading) {
    return (
      <section className="page compare">
        <h1>Comparing {repo1} vs {repo2}…</h1>
        <div className="spinner" />
      </section>
    )
  }

  if (error) {
    return (
      <section className="page compare">
        <h1>Compare</h1>
        <p className="error-msg">{error}</p>
      </section>
    )
  }

  if (!result) return null

  return (
    <section className="page compare">
      <h1>{repo1} vs {repo2}</h1>
      <CompareView result={result} history1={history1} history2={history2} />
    </section>
  )
}

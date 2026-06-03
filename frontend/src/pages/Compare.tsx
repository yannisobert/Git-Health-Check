import { useParams } from 'react-router-dom'
import { CompareView } from '../components/CompareView'

export function Compare() {
  const { o1 = '', r1 = '', o2 = '', r2 = '' } = useParams()

  const placeholder = {
    repo1: {
      owner: o1,
      repo: r1,
      fullName: `${o1}/${r1}`,
      score: 0,
      maxScore: 100,
      categories: [],
      analyzedAt: new Date().toISOString(),
    },
    repo2: {
      owner: o2,
      repo: r2,
      fullName: `${o2}/${r2}`,
      score: 0,
      maxScore: 100,
      categories: [],
      analyzedAt: new Date().toISOString(),
    },
    diff: {},
  }

  return (
    <section className="page compare">
      <h1>
        Compare {o1}/{r1} vs {o2}/{r2}
      </h1>
      <CompareView result={placeholder} />
    </section>
  )
}

import { useParams } from 'react-router-dom'
import { CheckList } from '../components/CheckList'
import { EvolutionChart } from '../components/EvolutionChart'
import { RivalSuggestions } from '../components/RivalSuggestions'
import { ScoreCard } from '../components/ScoreCard'

export function Report() {
  const { owner = '', repo = '' } = useParams()
  const fullName = `${owner}/${repo}`

  return (
    <section className="page report">
      <h1>Report — {fullName}</h1>
      <p className="placeholder">
        Full analysis will be available after feature/analyzer-core and feature/rest-server.
      </p>
      <ScoreCard repo={fullName} score={0} maxScore={100} />
      <CheckList checks={[]} />
      <EvolutionChart points={[]} />
      <RivalSuggestions rivals={[]} />
    </section>
  )
}

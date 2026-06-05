import { Link, useParams } from 'react-router-dom'
import type { RivalSuggestion } from '../api/types'

interface RivalSuggestionsProps {
  rivals: RivalSuggestion[]
}

export function RivalSuggestions({ rivals }: RivalSuggestionsProps) {
  const { owner = '', repo = '' } = useParams()

  const filtered = rivals.filter((r) => r.fullName)
  if (filtered.length === 0) return null

  return (
    <section className="rivals">
      <h3>Compare with similar repos</h3>
      <ul className="rival-list">
        {filtered.map((rival) => {
          const [rOwner, rRepo] = rival.fullName.split('/')
          return (
            <li key={rival.fullName} className="rival-item">
              <Link to={`/compare/${owner}/${repo}/${rOwner}/${rRepo}`}>
                {rival.fullName}
              </Link>
              <span className="rival-reason">{rival.reason}</span>
            </li>
          )
        })}
      </ul>
    </section>
  )
}

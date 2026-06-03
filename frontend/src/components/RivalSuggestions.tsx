import { Link } from 'react-router-dom'
import type { RivalSuggestion } from '../api/types'

interface RivalSuggestionsProps {
  rivals: RivalSuggestion[]
}

export function RivalSuggestions({ rivals }: RivalSuggestionsProps) {
  if (rivals.length === 0) return null

  return (
    <section className="rivals">
      <h3>Compare with similar repos</h3>
      <ul>
        {rivals.map((rival) => {
          const [owner, repo] = rival.fullName.split('/')
          return (
            <li key={rival.fullName}>
              <Link to={`/compare/${owner}/${repo}/placeholder/placeholder`}>
                {rival.fullName}
              </Link>
              <span> — {rival.reason}</span>
            </li>
          )
        })}
      </ul>
    </section>
  )
}

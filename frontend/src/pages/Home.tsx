import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const EXAMPLES = [
  'gin-gonic/gin',
  'facebook/react',
  'django/django',
  'expressjs/express',
]

export function Home() {
  const [url, setUrl] = useState('')
  const navigate = useNavigate()

  function parseRepo(input: string): string | null {
    const trimmed = input.trim()
    const match = trimmed.match(/github\.com\/([^/]+\/[^/\s#?]+)/i)
    if (match) return match[1].replace(/\.git$/, '')
    if (/^[^/]+\/[^/]+$/.test(trimmed)) return trimmed
    return null
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const repo = parseRepo(url)
    if (!repo) return
    const [owner, name] = repo.split('/')
    navigate(`/report/${owner}/${name}`)
  }

  return (
    <section className="page home">
      <h1>GitHub Repository Health Check</h1>
      <p>Analyze public repos and get a score out of 100 with actionable suggestions.</p>

      <form onSubmit={handleSubmit} className="repo-form">
        <input
          type="text"
          placeholder="https://github.com/owner/repo or owner/repo"
          value={url}
          onChange={(e) => setUrl(e.target.value)}
        />
        <button type="submit">Analyze</button>
      </form>

      <div className="examples">
        <p>Try an example:</p>
        {EXAMPLES.map((repo) => {
          const [owner, name] = repo.split('/')
          return (
            <button
              key={repo}
              type="button"
              className="link-button"
              onClick={() => navigate(`/report/${owner}/${name}`)}
            >
              {repo}
            </button>
          )
        })}
      </div>
    </section>
  )
}

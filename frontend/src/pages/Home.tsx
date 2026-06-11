import { useState } from 'react'
import { useNavigate } from 'react-router-dom'

const EXAMPLES = ['gin-gonic/gin', 'facebook/react', 'django/django', 'expressjs/express']

const CATEGORIES = [
  {
    icon: '🔧',
    name: 'Maintenance',
    pts: 30,
    desc: 'Open issues, PR merge time, last release, dependency freshness.',
  },
  {
    icon: '⚡',
    name: 'Activity',
    pts: 25,
    desc: 'Commit frequency, contributor count, recent pushes.',
  },
  {
    icon: '📐',
    name: 'Conventions',
    pts: 25,
    desc: 'README, license, changelog, contributing guide, semantic versioning.',
  },
  {
    icon: '🚀',
    name: 'CI / CD',
    pts: 20,
    desc: 'GitHub Actions workflows, branch protection, automated checks.',
  },
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
    <div>
      <section className="hero">
        <h1>Analyze any GitHub repo</h1>
        <p className="hero-sub">
          Get a health score out of 100 — maintenance, activity, conventions, CI/CD.
        </p>
        <form onSubmit={handleSubmit} className="search-form">
          <input
            type="text"
            placeholder="owner/repo or https://github.com/owner/repo"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            autoFocus
          />
          <button type="submit" className="search-btn">Analyze →</button>
        </form>
        <div className="example-chips">
          <span>Try:</span>
          {EXAMPLES.map((repo) => {
            const [owner, name] = repo.split('/')
            return (
              <button
                key={repo}
                type="button"
                className="example-chip"
                onClick={() => navigate(`/report/${owner}/${name}`)}
              >
                {repo}
              </button>
            )
          })}
        </div>
      </section>

      <div className="home-stats">
        <span>100 point scale</span>
        <span className="home-stats-dot" />
        <span>4 categories</span>
        <span className="home-stats-dot" />
        <span>Public repos</span>
        <span className="home-stats-dot" />
        <span>No sign-up</span>
      </div>

      <div className="home-features">
        <h2 className="home-features-title">What we check</h2>
        <div className="features-grid">
          {CATEGORIES.map((cat) => (
            <div key={cat.name} className="feature-card">
              <div className="feature-icon">{cat.icon}</div>
              <div className="feature-body">
                <div className="feature-header">
                  <strong className="feature-name">{cat.name}</strong>
                  <span className="feature-pts">{cat.pts} pts</span>
                </div>
                <p className="feature-desc">{cat.desc}</p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

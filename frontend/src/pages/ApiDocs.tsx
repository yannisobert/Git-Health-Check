const ENDPOINTS = [
  {
    method: 'GET',
    path: '/api/analyze?repo=owner/repo',
    desc: 'Full health report — score, checks, and improvement suggestions.',
  },
  {
    method: 'GET',
    path: '/api/history?repo=owner/repo&period=weekly',
    desc: 'Score evolution over time. period: weekly · monthly · 2year · 3year',
  },
  {
    method: 'GET',
    path: '/api/compare?repo1=a/b&repo2=c/d',
    desc: 'Side-by-side comparison with score diff per category.',
  },
  {
    method: 'GET',
    path: '/api/rivals?repo=owner/repo',
    desc: 'Suggested similar repositories to compare against.',
  },
  {
    method: 'GET',
    path: '/api/health',
    desc: 'Server health check.',
  },
]

export function ApiDocs() {
  return (
    <>
      <div className="page-banner">
        <div className="page-banner-inner">
          <p className="page-banner-eyebrow">Developer</p>
          <h1>REST API</h1>
          <p className="banner-sub">
            All endpoints return JSON. Public repos require no authentication.
          </p>
        </div>
      </div>

      <div className="container-overlap">
        <div style={{ marginBottom: '1.5rem' }}>
          <p style={{ fontSize: '0.8rem', color: 'var(--text-muted)', margin: '0 0 0.25rem' }}>
            Add a <code style={{ fontWeight: 600, color: '#6366f1' }}>GITHUB_TOKEN</code> environment variable to raise the rate limit from{' '}
            <strong>60</strong> to <strong>5 000</strong> req/h.
          </p>
        </div>

        {ENDPOINTS.map((ep) => (
          <div key={ep.path} className="endpoint-card">
            <span className={`method-badge method-${ep.method.toLowerCase()}`}>{ep.method}</span>
            <div>
              <p className="endpoint-path">{ep.path}</p>
              <p className="endpoint-desc">{ep.desc}</p>
            </div>
          </div>
        ))}
      </div>
    </>
  )
}

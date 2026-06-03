const ENDPOINTS = [
  { method: 'GET', path: '/api/analyze?repo=owner/repo', desc: 'Full health report' },
  { method: 'GET', path: '/api/history?repo=owner/repo&period=weekly', desc: 'Score history' },
  { method: 'GET', path: '/api/compare?repo1=a/b&repo2=c/d', desc: 'Compare two repos' },
  { method: 'GET', path: '/api/rivals?repo=owner/repo', desc: 'Similar repo suggestions' },
  { method: 'GET', path: '/api/health', desc: 'Server health check' },
]

export function ApiDocs() {
  return (
    <section className="page api-docs">
      <h1>REST API</h1>
      <ul className="endpoint-list">
        {ENDPOINTS.map((ep) => (
          <li key={ep.path}>
            <code>
              {ep.method} {ep.path}
            </code>
            <span>{ep.desc}</span>
          </li>
        ))}
      </ul>
    </section>
  )
}

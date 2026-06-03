import { Link } from 'react-router-dom'

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="app">
      <header className="header">
        <Link to="/" className="logo">
          ghhealth
        </Link>
        <nav className="nav">
          <Link to="/">Home</Link>
          <Link to="/history">History</Link>
          <Link to="/settings">Settings</Link>
          <Link to="/api">API</Link>
        </nav>
      </header>
      <main className="main">{children}</main>
    </div>
  )
}

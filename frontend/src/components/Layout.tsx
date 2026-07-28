import { Link } from 'react-router-dom'
import logo from '../assets/logo.png'
import { Footer } from './Footer'

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="app">
      <header className="navbar">
        <div className="navbar-inner">
          <Link to="/" className="logo">
            <img src={logo} alt="" height={36} />
            <span>Owlspector</span>
          </Link>
          <nav className="nav">
            <Link to="/">Home</Link>
            <Link to="/settings">Settings</Link>
            <Link to="/api">API</Link>
          </nav>
        </div>
      </header>
      <div className="page-wrapper">
        {children}
      </div>
      <Footer />
    </div>
  )
}

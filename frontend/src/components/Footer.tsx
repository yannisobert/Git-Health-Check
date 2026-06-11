import logo from '../assets/logo.png'

export function Footer() {
  return (
    <footer className="footer">
      <div className="footer-inner">
        <div className="footer-brand">
          <img src={logo} alt="ghhealth" height={22} />
          <span className="footer-copy">© {new Date().getFullYear()} ghhealth</span>
        </div>
        <div className="footer-links">
          <a href="https://github.com/yannisobert/Git-Health-Check" target="_blank" rel="noreferrer">
            GitHub
          </a>
          <a href="/api">API</a>
        </div>
      </div>
    </footer>
  )
}

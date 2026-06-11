export function Settings() {
  return (
    <>
      <div className="page-banner">
        <div className="page-banner-inner">
          <p className="page-banner-eyebrow">Configuration</p>
          <h1>Settings</h1>
          <p className="banner-sub">Customize your analysis experience.</p>
        </div>
      </div>

      <div className="container-overlap">
        <div className="settings-section">
          <p className="settings-section-title">Access</p>
          <form className="settings-form">
            <label>
              GitHub token
              <input type="password" placeholder="ghp_…" disabled />
              <span className="hint">Raises the rate limit from 60 to 5 000 req/h. Free to generate on GitHub.</span>
            </label>
          </form>
        </div>

        <div className="settings-section">
          <p className="settings-section-title">Analysis</p>
          <form className="settings-form">
            <label>
              Default period
              <select disabled>
                <option value="weekly">Weekly (26 weeks)</option>
                <option value="monthly">Monthly (1 year)</option>
                <option value="2year">2 Years</option>
                <option value="3year">3 Years</option>
              </select>
            </label>
          </form>
        </div>

        <div className="settings-notice">
          <span>💡</span>
          <span>Settings persistence is coming in a later phase — values are not saved yet.</span>
        </div>
      </div>
    </>
  )
}

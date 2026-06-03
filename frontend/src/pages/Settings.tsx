export function Settings() {
  return (
    <section className="page settings">
      <h1>Settings</h1>
      <form className="settings-form">
        <label>
          GitHub token (optional)
          <input type="password" placeholder="ghp_..." disabled />
        </label>
        <label>
          Analysis period
          <select disabled>
            <option value="weekly">Weekly</option>
            <option value="monthly">Monthly</option>
          </select>
        </label>
        <p className="hint">Settings persistence will be added in a later phase.</p>
      </form>
    </section>
  )
}

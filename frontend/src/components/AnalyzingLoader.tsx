import { useEffect, useState } from 'react'
import logo from '../assets/logo.png'

const STEPS = [
  'Fetching repository info',
  'Analyzing commit history',
  'Checking issues & pull requests',
  'Scanning CI/CD workflows',
  'Evaluating conventions',
  'Computing score',
]

export function AnalyzingLoader({ repo }: { repo: string }) {
  const [activeStep, setActiveStep] = useState(0)

  useEffect(() => {
    const id = setInterval(() => {
      setActiveStep((s) => (s < STEPS.length - 1 ? s + 1 : s))
    }, 700)
    return () => clearInterval(id)
  }, [])

  return (
    <div className="analyzing-page">
      <div className="analyzing-ring-wrap">
        <div className="analyzing-ring">
          <svg viewBox="0 0 100 100" width="96" height="96">
            <defs>
              <linearGradient id="analyzing-grad" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stopColor="#6366f1" />
                <stop offset="100%" stopColor="#a855f7" />
              </linearGradient>
            </defs>
            <circle cx="50" cy="50" r="42" fill="none" stroke="#f1f5f9" strokeWidth="6" />
            <circle
              cx="50" cy="50" r="42"
              fill="none"
              stroke="url(#analyzing-grad)"
              strokeWidth="6"
              strokeLinecap="round"
              strokeDasharray="198 66"
            />
          </svg>
        </div>
        <img src={logo} alt="" className="analyzing-logo" />
      </div>

      <p className="analyzing-repo">{repo}</p>

      <div className="analyzing-steps">
        {STEPS.map((step, i) => (
          <div
            key={step}
            className={`analyzing-step ${
              i < activeStep ? 'step-done' : i === activeStep ? 'step-active' : 'step-pending'
            }`}
          >
            <span className="step-dot" />
            <span>{step}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

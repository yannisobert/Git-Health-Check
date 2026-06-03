import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { Layout } from './components/Layout'
import { ApiDocs } from './pages/ApiDocs'
import { Compare } from './pages/Compare'
import { History } from './pages/History'
import { Home } from './pages/Home'
import { Report } from './pages/Report'
import { Settings } from './pages/Settings'

export default function App() {
  return (
    <BrowserRouter>
      <Layout>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/report/:owner/:repo" element={<Report />} />
          <Route path="/compare/:o1/:r1/:o2/:r2" element={<Compare />} />
          <Route path="/history" element={<History />} />
          <Route path="/settings" element={<Settings />} />
          <Route path="/api" element={<ApiDocs />} />
        </Routes>
      </Layout>
    </BrowserRouter>
  )
}

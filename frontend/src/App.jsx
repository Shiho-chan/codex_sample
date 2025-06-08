import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import Home from './pages/Home.jsx'
import Board from './pages/Board.jsx'
import './App.css'

function App() {
  return (
    <Router>
      <nav>
        <Link to="/">Home</Link> | <Link to="/board">掲示板</Link>
      </nav>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/board" element={<Board />} />
      </Routes>
    </Router>
  )
}

export default App

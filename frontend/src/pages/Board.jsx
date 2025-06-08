import { useEffect, useState } from 'react'
import '../App.css'

function Board() {
  const [comments, setComments] = useState([])
  const [text, setText] = useState('')
  const [name, setName] = useState('')
  const [user, setUser] = useState(null)
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')

  const fetchComments = async () => {
    try {
      const res = await fetch('/api/comments')
      if (!res.ok) return
      setComments(await res.json())
    } catch {
      // ignore
    }
  }

  useEffect(() => {
    fetch('/api/me')
      .then(async res => {
        if (res.ok) setUser(await res.json())
      })
      .catch(() => {})
    fetchComments()
  }, [])

  const login = async (e) => {
    e.preventDefault()
    try {
      const res = await fetch('/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password })
      })
      if (res.ok) setUser(await res.json())
    } catch {
      // ignore
    }
  }

  const submit = async (e) => {
    e.preventDefault()
    try {
      await fetch('/api/comments', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, text })
      })
    } catch {
      // ignore
    }
    setText('')
    setName('')
    fetchComments()
  }

  const remove = async (id) => {
    try {
      await fetch(`/api/comments/${id}`, { method: 'DELETE' })
      fetchComments()
    } catch {
      // ignore
    }
  }

  return (
    <div className="board">
      <h1>掲示板</h1>
      {user ? (
        <p>Logged in as {user.username} ({user.role})</p>
      ) : (
        <form onSubmit={login}>
          <input placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} />
          <input type="password" placeholder="Password" value={password} onChange={e => setPassword(e.target.value)} />
          <button type="submit">Login</button>
        </form>
      )}
      <form onSubmit={submit}>
        <input placeholder="Handle" value={name} onChange={e => setName(e.target.value)} />
        <input value={text} onChange={e => setText(e.target.value)} />
        <button type="submit">投稿</button>
      </form>
      <ul>
        {comments.map(c => (
          <li key={c.id}>
            <span>{c.name} ({new Date(c.timestamp).toLocaleString()}): {c.text}</span>
            {user && (user.role === 'admin' || user.username === c.username) && (
              <button onClick={() => remove(c.id)} className="delete">削除</button>
            )}
          </li>
        ))}
      </ul>
    </div>
  )
}

export default Board

import { useEffect, useState } from 'react'
import '../App.css'

function Board() {
  const [comments, setComments] = useState([])
  const [text, setText] = useState('')

  useEffect(() => {
    fetch('/api/comments')
      .then(res => res.json())
      .then(setComments)
  }, [])

  const submit = async (e) => {
    e.preventDefault()
    await fetch('/api/comments', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text })
    })
    setText('')
    const res = await fetch('/api/comments')
    setComments(await res.json())
  }

  const remove = async (id) => {
    await fetch(`/api/comments/${id}`, { method: 'DELETE' })
    const res = await fetch('/api/comments')
    setComments(await res.json())
  }

  return (
    <div className="board">
      <h1>掲示板</h1>
      <form onSubmit={submit}>
        <input value={text} onChange={e => setText(e.target.value)} />
        <button type="submit">投稿</button>
      </form>
      <ul>
        {comments.map(c => (
          <li key={c.id}>
            {c.text}
            <button onClick={() => remove(c.id)} className="delete">削除</button>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default Board

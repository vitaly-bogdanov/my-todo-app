import { useState } from 'react'

function filled(body) {
  return Object.fromEntries(Object.entries(body).filter(([, value]) => value !== '' && value !== undefined))
}

export function TasksView({ tasks, total, error, loading, onLoad, onCreate, onGet, onPatch, onRemove }) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [authorId, setAuthorId] = useState('')
  const [userId, setUserId] = useState('')
  const [limit, setLimit] = useState('')
  const [offset, setOffset] = useState('')
  const [id, setId] = useState('')
  const [patchTitle, setPatchTitle] = useState('')
  const [patchDescription, setPatchDescription] = useState('')
  const [completed, setCompleted] = useState('')

  return (
    <section>
      <h2>Tasks</h2>
      <form onSubmit={(event) => {
        event.preventDefault()
        onCreate({
          title,
          description: description || undefined,
          author_user_id: Number(authorId),
        })
      }}>
        <input value={title} onChange={(event) => setTitle(event.target.value)} placeholder="Title" required />
        <input value={description} onChange={(event) => setDescription(event.target.value)} placeholder="Description" />
        <input value={authorId} onChange={(event) => setAuthorId(event.target.value)} placeholder="Author user id" required />
        <button disabled={loading}>Add task</button>
      </form>
      <form onSubmit={(event) => {
        event.preventDefault()
        onLoad({ user_id: userId, limit, offset })
      }}>
        <input value={userId} onChange={(event) => setUserId(event.target.value)} placeholder="Author user id" />
        <input value={limit} onChange={(event) => setLimit(event.target.value)} placeholder="Limit" />
        <input value={offset} onChange={(event) => setOffset(event.target.value)} placeholder="Offset" />
        <button disabled={loading}>Load tasks</button>
      </form>
      <form onSubmit={(event) => {
        event.preventDefault()
        onGet(id)
      }}>
        <input value={id} onChange={(event) => setId(event.target.value)} placeholder="Task id" required />
        <button disabled={loading}>Get task</button>
        <button type="button" disabled={loading} onClick={() => onRemove(id)}>Delete task</button>
      </form>
      <form onSubmit={(event) => {
        event.preventDefault()
        onPatch(id, filled({
          title: patchTitle,
          description: patchDescription,
          completed: completed === '' ? undefined : completed === 'true',
        }))
      }}>
        <input value={patchTitle} onChange={(event) => setPatchTitle(event.target.value)} placeholder="New title" />
        <input value={patchDescription} onChange={(event) => setPatchDescription(event.target.value)} placeholder="New description" />
        <select value={completed} onChange={(event) => setCompleted(event.target.value)}>
          <option value="">Completed unchanged</option>
          <option value="true">Completed</option>
          <option value="false">Not completed</option>
        </select>
        <button disabled={loading}>Patch task</button>
      </form>
      {total !== null && <p className="total">Total: {total}</p>}
      <ul id="tasks-list">
        {tasks.map((task) => (
          <li key={task.id}>
            {task.id}. {task.title}{task.completed ? ' (done)' : ''}
          </li>
        ))}
      </ul>
      <p className="message">{error}</p>
    </section>
  )
}

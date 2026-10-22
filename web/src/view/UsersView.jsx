import { useState } from 'react'

function filled(body) {
  return Object.fromEntries(Object.entries(body).filter(([, value]) => value !== ''))
}

export function UsersView({ users, total, error, loading, onLoad, onCreate, onGet, onPatch, onRemove }) {
  const [fullName, setFullName] = useState('')
  const [phone, setPhone] = useState('')
  const [limit, setLimit] = useState('')
  const [offset, setOffset] = useState('')
  const [id, setId] = useState('')
  const [patchName, setPatchName] = useState('')
  const [patchPhone, setPatchPhone] = useState('')

  return (
    <section>
      <h2>Users</h2>
      <form onSubmit={(event) => {
        event.preventDefault()
        onCreate(filled({ full_name: fullName, phone_number: phone }))
      }}>
        <input value={fullName} onChange={(event) => setFullName(event.target.value)} placeholder="Full name" required />
        <input value={phone} onChange={(event) => setPhone(event.target.value)} placeholder="Phone" />
        <button disabled={loading}>Add user</button>
      </form>
      <form onSubmit={(event) => {
        event.preventDefault()
        onLoad({ limit, offset })
      }}>
        <input value={limit} onChange={(event) => setLimit(event.target.value)} placeholder="Limit" />
        <input value={offset} onChange={(event) => setOffset(event.target.value)} placeholder="Offset" />
        <button disabled={loading}>Load users</button>
      </form>
      <form onSubmit={(event) => {
        event.preventDefault()
        onGet(id)
      }}>
        <input value={id} onChange={(event) => setId(event.target.value)} placeholder="User id" required />
        <button disabled={loading}>Get user</button>
        <button type="button" disabled={loading} onClick={() => onRemove(id)}>Delete user</button>
      </form>
      <form onSubmit={(event) => {
        event.preventDefault()
        onPatch(id, filled({ full_name: patchName, phone_number: patchPhone }))
      }}>
        <input value={patchName} onChange={(event) => setPatchName(event.target.value)} placeholder="New full name" />
        <input value={patchPhone} onChange={(event) => setPatchPhone(event.target.value)} placeholder="New phone" />
        <button disabled={loading}>Patch user</button>
      </form>
      {total !== null && <p className="total">Total: {total}</p>}
      <ul id="users-list">
        {users.map((user) => (
          <li key={user.id}>
            {user.id}. {user.full_name}{user.phone_number ? `, ${user.phone_number}` : ''}
          </li>
        ))}
      </ul>
      <p className="message">{error}</p>
    </section>
  )
}

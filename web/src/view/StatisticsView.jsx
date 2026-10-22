import { useState } from 'react'

export function StatisticsView({ statistics, error, loading, onLoad }) {
  const [userId, setUserId] = useState('')
  const [from, setFrom] = useState('')
  const [to, setTo] = useState('')

  return (
    <section>
      <h2>Statistics</h2>
      <form onSubmit={(event) => {
        event.preventDefault()
        onLoad({ user_id: userId, from, to })
      }}>
        <input value={userId} onChange={(event) => setUserId(event.target.value)} placeholder="Author user id" />
        <input value={from} onChange={(event) => setFrom(event.target.value)} placeholder="From YYYY-MM-DD" />
        <input value={to} onChange={(event) => setTo(event.target.value)} placeholder="To YYYY-MM-DD" />
        <button disabled={loading}>Load statistics</button>
      </form>
      <ul id="statistics-list">
        {statistics && (
          <>
            <li>Created: {statistics.tasks_created}</li>
            <li>Completed: {statistics.tasks_completed}</li>
            <li>Rate: {statistics.tasks_completed_rate ?? '—'}</li>
            <li>Average time: {statistics.tasks_average_completation_time ?? '—'}</li>
          </>
        )}
      </ul>
      <p className="message">{error}</p>
    </section>
  )
}

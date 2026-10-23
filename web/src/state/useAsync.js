import { useCallback, useState } from 'react'

export function useAsync() {
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  const run = useCallback(async (action) => {
    setLoading(true)
    setError('')
    try {
      return await action()
    } catch (err) {
      setError(err.message)
      return null
    } finally {
      setLoading(false)
    }
  }, [])

  return { error, loading, run }
}

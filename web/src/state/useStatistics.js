import { useCallback, useState } from 'react'
import * as statisticsApi from '../api/statistics'
import { useAsync } from './useAsync'

export function useStatistics() {
  const [statistics, setStatistics] = useState(null)
  const { error, loading, run } = useAsync()

  const load = useCallback((params) => run(async () => {
    const result = await statisticsApi.getStatistics(params)
    setStatistics(result.data)
  }), [run])

  return { statistics, error, loading, load }
}

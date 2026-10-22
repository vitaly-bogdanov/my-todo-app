import { useCallback, useState } from 'react'
import * as tasksApi from '../api/tasks'
import { useAsync } from './useAsync'

export function useTasks() {
  const [tasks, setTasks] = useState([])
  const [total, setTotal] = useState(null)
  const { error, loading, run } = useAsync()

  const load = useCallback((params) => run(async () => {
    const result = await tasksApi.listTasks(params)
    setTasks(result.data ?? [])
    setTotal(result.total)
  }), [run])

  const create = useCallback((body) => run(async () => {
    await tasksApi.createTask(body)
    const result = await tasksApi.listTasks()
    setTasks(result.data ?? [])
    setTotal(result.total)
  }), [run])

  const getOne = useCallback((id) => run(async () => {
    const result = await tasksApi.getTask(id)
    setTasks(result.data ? [result.data] : [])
    setTotal(null)
  }), [run])

  const patch = useCallback((id, body) => run(async () => {
    const result = await tasksApi.patchTask(id, body)
    setTasks(result.data ? [result.data] : [])
    setTotal(null)
  }), [run])

  const remove = useCallback((id) => run(async () => {
    await tasksApi.deleteTask(id)
    const result = await tasksApi.listTasks()
    setTasks(result.data ?? [])
    setTotal(result.total)
  }), [run])

  return { tasks, total, error, loading, load, create, getOne, patch, remove }
}

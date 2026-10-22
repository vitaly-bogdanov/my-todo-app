import { useCallback, useState } from 'react'
import * as usersApi from '../api/users'
import { useAsync } from './useAsync'

export function useUsers() {
  const [users, setUsers] = useState([])
  const [total, setTotal] = useState(null)
  const { error, loading, run } = useAsync()

  const load = useCallback((params) => run(async () => {
    const result = await usersApi.listUsers(params)
    setUsers(result.data ?? [])
    setTotal(result.total)
  }), [run])

  const create = useCallback((body) => run(async () => {
    await usersApi.createUser(body)
    const result = await usersApi.listUsers()
    setUsers(result.data ?? [])
    setTotal(result.total)
  }), [run])

  const getOne = useCallback((id) => run(async () => {
    const result = await usersApi.getUser(id)
    setUsers(result.data ? [result.data] : [])
    setTotal(null)
  }), [run])

  const patch = useCallback((id, body) => run(async () => {
    const result = await usersApi.patchUser(id, body)
    setUsers(result.data ? [result.data] : [])
    setTotal(null)
  }), [run])

  const remove = useCallback((id) => run(async () => {
    await usersApi.deleteUser(id)
    const result = await usersApi.listUsers()
    setUsers(result.data ?? [])
    setTotal(result.total)
  }), [run])

  return { users, total, error, loading, load, create, getOne, patch, remove }
}

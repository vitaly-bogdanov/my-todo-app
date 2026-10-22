import { request, withQuery } from './http'

export function listUsers(params) {
  return request(withQuery('/users', params))
}

export function getUser(id) {
  return request(`/users/${id}`)
}

export function createUser(body) {
  return request('/users', { method: 'POST', body })
}

export function patchUser(id, body) {
  return request(`/users/${id}`, { method: 'PATCH', body })
}

export function deleteUser(id) {
  return request(`/users/${id}`, { method: 'DELETE' })
}

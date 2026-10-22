import { request, withQuery } from './http'

export function listTasks(params) {
  return request(withQuery('/tasks', params))
}

export function getTask(id) {
  return request(`/tasks/${id}`)
}

export function createTask(body) {
  return request('/tasks', { method: 'POST', body })
}

export function patchTask(id, body) {
  return request(`/tasks/${id}`, { method: 'PATCH', body })
}

export function deleteTask(id) {
  return request(`/tasks/${id}`, { method: 'DELETE' })
}

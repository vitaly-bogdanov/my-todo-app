const baseURL = '/api/v1'

export function withQuery(path, params = {}) {
  const query = new URLSearchParams()
  for (const [key, value] of Object.entries(params)) {
    if (value !== '' && value !== null && value !== undefined) {
      query.set(key, String(value))
    }
  }
  const text = query.toString()
  return text ? `${path}?${text}` : path
}

export async function request(path, { method = 'GET', body } = {}) {
  const response = await fetch(`${baseURL}${path}`, {
    method,
    headers: body === undefined ? undefined : { 'Content-Type': 'application/json' },
    body: body === undefined ? undefined : JSON.stringify(body),
  })

  const text = await response.text()
  const data = text ? JSON.parse(text) : null
  if (!response.ok) {
    throw new Error(data?.message || data?.error || response.statusText)
  }

  return {
    data,
    total: response.headers.get('X-Total-Count'),
  }
}

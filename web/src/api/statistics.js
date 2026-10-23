import { request, withQuery } from './http'

export function getStatistics(params) {
  return request(withQuery('/statistics', params))
}

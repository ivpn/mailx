import { api } from './api'

export const subscriptionApi = {
    get: () => api.get('/sub'),
}
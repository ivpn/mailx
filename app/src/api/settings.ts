import { api } from './api'

export const settingsApi = {
    get: () => api.get('/settings'),
    update: (data: any) => api.put('/settings', data),
}
import { api } from './api'

export const settingsApi = {
    get: () => api.get('/settings'),
    getDefaults: () => api.get('/defaults'),
    update: (data: any) => api.put('/settings', data),
}
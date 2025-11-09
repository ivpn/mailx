import { api } from './api'

export const discardApi = {
    getList: () => api.get('/discards'),
}
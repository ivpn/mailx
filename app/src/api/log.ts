import { api } from './api'

export const logApi = {
    getList: () => api.get('/logs'),
    getFile: (id: string) => api.get(`/log/file/${id}`),
    deleteAll: () => api.delete('/logs'),
}
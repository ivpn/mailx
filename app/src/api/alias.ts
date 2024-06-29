import { api } from './api'

export const aliasApi = {
    get: (id: string) => api.get('/alias/' + id),
    getList: (data: any) => api.get('/aliases', { params: data }),
    create: (data: any) => api.post('/alias', data),
    update: (id: string, data: any) => api.put('/alias/' + id, data),
    delete: (id: string) => api.delete('/alias/' + id),
}
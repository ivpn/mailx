import { api } from './api'

export const blocklistApi = {
    getList: (data: any) => api.get('/blocklist', { params: data }),
    create: (data: any) => api.post('/blocklist', data),
    delete: (id: string) => api.delete('/blocklist/' + id),
}

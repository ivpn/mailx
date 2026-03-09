import { api } from './api'

export const domainApi = {
    getList: () => api.get('/domains'),
    getConfig: () => api.get('/domains/dns-config'),
    create: (data: any) => api.post('/domain', data),
    update: (id: string, data: any) => api.put('/domain/' + id, data),
    delete: (id: string) => api.put('/domain/' + id),
}
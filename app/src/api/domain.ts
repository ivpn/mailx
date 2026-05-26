import { api } from './api'

export const domainApi = {
    getList: () => api.get('/domains'),
    getConfig: () => api.get('/domains/dns-config'),
    verifyDns: (id: string) => api.post('/domain/' + id + '/verify-dns'),
    create: (data: any) => api.post('/domain', data),
    update: (id: string, data: any) => api.put('/domain/' + id, data),
    delete: (id: string) => api.delete('/domain/' + id),
}
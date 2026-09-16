import { api } from './api'

export const aliasApi = {
    get: (id: string) => api.get('/alias/' + id),
    getList: (data: any) => api.get('/aliases', { params: data }),
    getWildcardDomainInfo: (domain: string) => api.get('/alias/wildcard-domain-info', { params: { domain } }),
    import: (data: any) => api.post('/aliases/import', data),
    export: () => api.get('/aliases/export'),
    create: (data: any) => api.post('/alias', data),
    update: (id: string, data: any) => api.put('/alias/' + id, data),
    delete: (id: string) => api.delete('/alias/' + id),
    restore: (id: string) => api.post('/alias/restore/' + id),
    forget: (id: string) => api.delete('/alias/forget/' + id),
    pin: (id: string, pinned: boolean) => api.put('/alias/' + id + '/pin', { pinned }),
}
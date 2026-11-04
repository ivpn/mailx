import { api } from './api'

export const ruleApi = {
    getList: () => api.get('/rules'),
    create: (data: any) => api.post('/rule', data),
    update: (id: string, data: any) => api.put('/rule/' + id, data),
    updateActive: (id: string, active: boolean) => api.put('/rule/' + id + '/active', { active }),
    reorder: (ids: string[]) => api.put('/rules/order', { ids }),
    delete: (id: string) => api.delete('/rule/' + id),
}

import { api } from './api'

export const recipientApi = {
    get: (id: string) => api.get('/recipient/' + id),
    getList: () => api.get('/recipients'),
    create: (data: any) => api.post('/recipient', data),
    update: (data: any) => api.put('/recipient', data),
    sendOtp: (id: string) => api.post('/recipient/sendotp/' + id),
    activate: (id: string, data: any) => api.post('/recipient/activate/' + id, data),
    delete: (id: string, data: any) => api.put('/recipient/delete/' + id, data),
}
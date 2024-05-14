import { api } from './api'

export const userApi = {
    register: (data: any) => api.post('/p/register', data),
    login: (data: any) => api.post('/p/login', data),
    logout: () => api.post('/user/logout'),
    delete: (data: any) => api.delete('/user/delete', data),
    sendOtp: () => api.post('/user/sendotp'),
    activate: (data: any) => api.post('/user/activate', data),
    get: () => api.get('/user'),
    stats: () => api.get('/user/stats'),
}
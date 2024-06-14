import { api } from './api'
import { setCookie } from 'typescript-cookie'

export const userApi = {
    register: (data: any) => api.post('/p/register', data),
    login: (data: any) => api.post('/p/login', data),
    logout: () => api.post('/user/logout'),
    delete: (data: any) => api.post('/user/delete', data),
    sendOtp: () => api.post('/user/sendotp'),
    activate: (data: any) => api.post('/user/activate', data),
    get: () => api.get('/user'),
    stats: () => api.get('/user/stats'),
    changePassword: (data: any) => api.put('/user/changepassword', data),
    initiatePasswordReset: (data: any) => api.post('/p/initiatepasswordreset', data),
    resetPassword: (data: any) => api.put('/p/resetpassword', data),
    clearSession: () => {
        setCookie('auth', '')
        localStorage.removeItem('email')
        window.location.href = '/login'
    },
}

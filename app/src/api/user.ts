import { api } from './api'

export const userApi = {
    register: (data: any) => api.post('/register', data),
    registerBegin: (data: any) => api.post('/register/begin', data),
    registerFinish: (data: any) => api.post('/register/finish', data),
    registerAdd: () => api.post('/register/add'),
    registerAddFinish: (data: any) => api.post('/register/add/finish', data),
    login: (data: any) => api.post('/login', data),
    loginBegin: (data: any) => api.post('/login/begin', data),
    loginFinish: (data: any) => api.post('/login/finish', data),
    logout: () => api.post('/user/logout'),
    deleteRequest: () => api.post('/user/delete/request'),
    delete: (data: any) => api.post('/user/delete', data),
    sendOtp: () => api.post('/user/sendotp'),
    activate: (data: any) => api.post('/user/activate', data),
    get: () => api.get('/user'),
    getCredentials: () => api.get('/user/credentials'),
    deleteCredential: (id: string) => api.delete('/user/credential/' + id),
    stats: () => api.get('/user/stats'),
    changePassword: (data: any) => api.put('/user/changepassword', data),
    changeEmail: (data: any) => api.put('/user/changeemail', data),
    initiatePasswordReset: (data: any) => api.post('/initiatepasswordreset', data),
    resetPassword: (data: any) => api.put('/resetpassword', data),
    totpEnable: () => api.put('/user/totp/enable'),
    totpEnableConfirm: (data: any) => api.put('/user/totp/enable/confirm', data),
    totpDisable: (data: any) => api.put('/user/totp/disable', data),
    accessKeyList: () => api.get('/accesskeys'),
    accessKeyCreate: (data: any) => api.post('/accesskeys', data),
    accessKeyDelete: (id: string) => api.delete('/accesskeys/' + id),
    clearSession: () => {
        localStorage.removeItem('email')
        window.location.href = '/login'
    },
}

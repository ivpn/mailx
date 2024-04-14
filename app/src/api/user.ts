import axios from 'axios'

let baseURL = 'http://localhost:3000/v1'

export const userApi = {
    register: (data: any) => axios.post(baseURL + '/register', data),
    login: (data: any) => axios.post(baseURL + '/login', data),
    logout: () => axios.post(baseURL + '/user/logout'),
    delete: (data: any) => axios.delete(baseURL + '/user/delete', data),
    sendOtp: () => axios.post(baseURL + '/user/sendotp'),
    activate: (data: any) => axios.post(baseURL + '/user/activate', data),
    stats: () => axios.get(baseURL + '/user/stats'),
}
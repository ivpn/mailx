import axios from 'axios'

let baseURL = 'http://localhost:3000'

export const userApi = {
    register: (data: any) => axios.post(baseURL + '/v1/register', data),
    login: (data: any) => axios.post(baseURL + '/v1/login', data),
    logout: () => axios.post(baseURL + '/v1/user/logout'),
    delete: (data: any) => axios.delete(baseURL + '/v1/user/delete', data),
    sendOtp: () => axios.post(baseURL + '/v1/user/sendotp'),
    activate: (data: any) => axios.post(baseURL + '/v1/user/activate', data),
    stats: () => axios.get(baseURL + '/v1/user/stats'),
}
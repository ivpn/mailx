import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL

export const userApi = {
    register: (data: any) => axios.post(apiUrl + '/v1/register', data),
    login: (data: any) => axios.post(apiUrl + '/v1/login', data),
    logout: () => axios.post(apiUrl + '/v1/user/logout'),
    delete: (data: any) => axios.delete(apiUrl + '/v1/user/delete', data),
    sendOtp: () => axios.post(apiUrl + '/v1/user/sendotp'),
    activate: (data: any) => axios.post(apiUrl + '/v1/user/activate', data),
    stats: () => axios.get(apiUrl + '/v1/user/stats'),
}
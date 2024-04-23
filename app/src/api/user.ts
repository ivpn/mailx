import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL + '/v1'

export const userApi = {
    register: (data: any) => axios.post(apiUrl + '/register', data),
    login: (data: any) => axios.post(apiUrl + '/login', data),
    logout: () => axios.post(apiUrl + '/user/logout', { withCredentials: true }),
    delete: (data: any) => axios.delete(apiUrl + '/user/delete', { ...data, withCredentials: true }),
    sendOtp: () => axios.post(apiUrl + '/user/sendotp', { withCredentials: true }),
    activate: (data: any) => axios.post(apiUrl + '/user/activate', { ...data, withCredentials: true }),
    get: () => axios.get(apiUrl + '/user', { withCredentials: true }),
    stats: () => axios.get(apiUrl + '/user/stats', { withCredentials: true }),
}
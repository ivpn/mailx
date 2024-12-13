import axios from 'axios'
import { userApi } from './user.ts'

export const api = axios.create({
    withCredentials: true,
    baseURL: import.meta.env.VITE_API_URL + '/v1'
})

api.interceptors.response.use(
    response => response, // simply return the response in case of success
    error => {
        if (error.response && error.response.status === 401 && window.location.pathname !== '/login') {
            // Handle the 401 error
            userApi.clearSession()
        }
        return Promise.reject(error)
    }
)

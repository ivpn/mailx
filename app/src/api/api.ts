import axios from 'axios'
import env from "../env.json"
import { setCookie } from 'typescript-cookie'

export const api = axios.create({
    withCredentials: true,
    baseURL: env.API_URL + '/v1'
})

api.interceptors.response.use(
    response => response, // simply return the response in case of success
    error => {
        if (error.response && error.response.status === 401) {
            // Handle the 401 error
            setCookie('auth', '')
            localStorage.removeItem('email')
            window.location.href = '/login'
        }
        return Promise.reject(error)
    }
)

import axios from 'axios'
import { userApi } from './user.ts'
import { requestStepUp } from '../stepup.ts'

export const api = axios.create({
    withCredentials: true,
    baseURL: import.meta.env.VITE_API_URL + '/v1'
})

api.interceptors.response.use(
    response => response, // simply return the response in case of success
    async error => {
        if (error.response && error.response.status === 401 && window.location.pathname.startsWith('/account')) {
            // Handle the 401 error
            userApi.clearSession()
            return Promise.reject(error)
        }

        if (error.response && error.response.data?.code === 80001 && !error.config._stepUpRetried) {
            try {
                await requestStepUp(error.response.data.methods)
                error.config._stepUpRetried = true
                return api(error.config)
            } catch {
                // Components read response.data.error, not the backend's original step-up message
                error.response.data.error = error.response.data.cancel_message
                return Promise.reject(error)
            }
        }

        return Promise.reject(error)
    }
)


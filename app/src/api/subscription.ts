import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL

export const subscriptionApi = {
    get: () => axios.get(apiUrl + '/v1/subscription', { withCredentials: true }),
}
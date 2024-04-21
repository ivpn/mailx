import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL + '/v1'

export const subscriptionApi = {
    get: () => axios.get(apiUrl + '/subscription', { withCredentials: true }),
}
import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL + '/v1'

export const settingsApi = {
    get: () => axios.get(apiUrl + '/settings', { withCredentials: true }),
    update: (data: any) => axios.put(apiUrl + '/settings', { ...data, withCredentials: true }),
}
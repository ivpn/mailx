import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL

export const settingsApi = {
    get: () => axios.get(apiUrl + '/settings'),
    update: (data: any) => axios.put(apiUrl + '/settings', data),
}
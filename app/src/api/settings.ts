import axios from 'axios'

const apiUrl = import.meta.env.VUE_APP_API_URL

export const settingsApi = {
    get: () => axios.get(apiUrl + '/settings'),
    update: (data: any) => axios.put(apiUrl + '/settings', data),
}
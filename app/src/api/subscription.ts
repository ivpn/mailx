import axios from 'axios'

const apiUrl = import.meta.env.VUE_APP_API_URL

export const subscriptionApi = {
    get: () => axios.get(apiUrl + '/subscription'),
}
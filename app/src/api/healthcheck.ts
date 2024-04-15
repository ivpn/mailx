import axios from 'axios'

const apiUrl = import.meta.env.VUE_APP_API_URL;

export const healthcheckApi = {
    livez: () => axios.get(apiUrl + '/livez'),
    readyz: () => axios.get(apiUrl + '/readyz'),
}
import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL

export const healthcheckApi = {
    livez: () => axios.get(apiUrl + '/livez'),
    readyz: () => axios.get(apiUrl + '/readyz'),
}
import axios from 'axios'

let baseURL = 'http://localhost:3000'

export const healthcheckApi = {
    livez: () => axios.get(baseURL + '/livez'),
    readyz: () => axios.get(baseURL + '/readyz'),
}
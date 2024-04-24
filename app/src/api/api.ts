import axios from 'axios'
import env from "../env.json"

export const api = axios.create({
    withCredentials: true,
    baseURL: env.API_URL + '/v1'
})
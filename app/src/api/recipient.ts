import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL + '/v1'

export const recipientApi = {
    get: (id: string) => axios.get(apiUrl + '/recipient/' + id, { withCredentials: true }),
    getList: () => axios.get(apiUrl + '/recipients', { withCredentials: true }),
    create: (data: any) => axios.post(apiUrl + '/recipient', { ...data, withCredentials: true }),
    sendOtp: (id: string) => axios.post(apiUrl + '/recipient/sendotp/' + id, { withCredentials: true }),
    activate: (id: string, data: any) => axios.post(apiUrl + '/recipient/activate/' + id, { ...data, withCredentials: true }),
    delete: (id: string) => axios.delete(apiUrl + '/recipient/' + id, { withCredentials: true }),
}
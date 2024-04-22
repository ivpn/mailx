import axios from 'axios'
import env from "../env.json"

const apiUrl = env.API_URL + '/v1'

export const aliasApi = {
    get: (id: string) => axios.get(apiUrl + '/alias/' + id, { withCredentials: true }),
    getList: () => axios.get(apiUrl + '/aliases', { withCredentials: true }),
    create: (data: any) => axios.post(apiUrl + '/alias', { ...data, withCredentials: true }),
    update: (id: string, data: any) => axios.put(apiUrl + '/alias/' + id, { ...data, withCredentials: true }),
    delete: (id: string) => axios.delete(apiUrl + '/alias/' + id, { withCredentials: true }),
}
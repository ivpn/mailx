import axios from 'axios'

const apiUrl = import.meta.env.VUE_APP_API_URL

export const aliasApi = {
    get: (id: string) => axios.get(apiUrl + '/alias/' + id),
    getList: () => axios.get(apiUrl + '/aliases'),
    create: (data: any) => axios.post(apiUrl + '/alias', data),
    update: (id: string, data: any) => axios.put(apiUrl + '/alias/' + id, data),
    delete: (id: string) => axios.delete(apiUrl + '/alias/' + id),
}
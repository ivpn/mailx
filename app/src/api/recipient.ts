import axios from 'axios'

const apiUrl = import.meta.env.VUE_APP_API_URL

export const recipientApi = {
    get: (id: string) => axios.get(apiUrl + '/recipient/' + id),
    getList: () => axios.get(apiUrl + '/recipients'),
    create: (data: any) => axios.post(apiUrl + '/recipient', data),
    sendOtp: (id: string) => axios.post(apiUrl + '/recipient/sendotp/' + id),
    activate: (id: string, data: any) => axios.post(apiUrl + '/recipient/activate/' + id, data),
    delete: (id: string) => axios.delete(apiUrl + '/recipient/' + id),
}
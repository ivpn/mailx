import { api } from './api'

export const bounceApi = {
    getList: () => api.get('/bounces'),
    getFile: (id: string) => api.get(`/bounce/file/${id}`),
}
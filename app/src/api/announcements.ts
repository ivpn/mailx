import { api } from './api.ts'

export const announcementsApi = {
    getList: () => api.get('/announcements'),
}

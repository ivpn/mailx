import { api } from './api'

export const subscriptionApi = {
    get: () => api.get('/sub'),
    update: (data: any) => api.put('/sub/update', data),
    rotateSessionId: (data: any) => api.put('/rotatepasession', data),
}
import { api } from './api'

export const healthcheckApi = {
    livez: () => api.get('/livez'),
    readyz: () => api.get('/readyz'),
}
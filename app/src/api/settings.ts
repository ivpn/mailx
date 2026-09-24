import { api } from './api'

export const settingsApi = {
    get: () => api.get('/settings'),
    getDefaults: () => api.get('/defaults'),
    update: (data: any) => api.put('/settings', data),
    updateListUnsubscribeBehaviour: (behaviour: number) =>
        api.put('/settings/list-unsubscribe-behaviour', { list_unsubscribe_behaviour: behaviour }),
}
import { api } from './api'

export const settingsApi = {
    get: () => api.get('/settings'),
    getDefaults: () => api.get('/defaults'),
    update: (data: any) => api.put('/settings', data),
    updateListUnsubscribeBehaviour: (behaviour: number) =>
        api.put('/settings/list-unsubscribe-behaviour', { list_unsubscribe_behaviour: behaviour }),
    updateBlockTrackers: (enabled: boolean) =>
        api.put('/settings/block-trackers', { block_trackers: enabled }),
    updateRemoveHeader: (removeHeader: boolean) =>
        api.put('/settings/remove-header', { remove_header: removeHeader }),
    updateLogIssues: (enabled: boolean) =>
        api.put('/settings/log-issues', { log_issues: enabled }),
}
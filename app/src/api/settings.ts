import { api } from './api'

export const settingsApi = {
    get: () => api.get('/settings'),
    getDefaults: () => api.get('/defaults'),
    updateDefaultDomain: (domain: string) =>
        api.put('/settings/domain', { domain }),
    updateDefaultRecipient: (recipient: string) =>
        api.put('/settings/recipient', { recipient }),
    updateAliasFormat: (format: string) =>
        api.put('/settings/alias-format', { alias_format: format }),
    updateFromName: (fromName: string) =>
        api.put('/settings/from-name', { from_name: fromName }),
    updateListUnsubscribeBehaviour: (behaviour: number) =>
        api.put('/settings/list-unsubscribe-behaviour', { list_unsubscribe_behaviour: behaviour }),
    updateBlockTrackers: (enabled: boolean) =>
        api.put('/settings/block-trackers', { block_trackers: enabled }),
    updateRemoveHeader: (removeHeader: boolean) =>
        api.put('/settings/remove-header', { remove_header: removeHeader }),
    updateLogIssues: (enabled: boolean) =>
        api.put('/settings/log-issues', { log_issues: enabled }),
}
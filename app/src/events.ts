import mitt from 'mitt'

// Define events
type Events = {
    'user.update': { email: string }
    'alias.create': {}
    'alias.update': {}
    'alias.enabled': { id: string, enabled: boolean }
    'alias.delete': { id: string, wildcard: boolean }
    'alias.forget': { id: string }
    'totp.enable': {}
    'totp.disable': {}
    'recipient.create': {}
    'recipient.update': {}
    'recipient.verify': {}
    'recipient.delete': { id: string }
    'recipient.delete.error': { error: string }
    'recipient.reload': {}
    'accesskey.create': {}
    'domain.create': {}
    'domain.reload': {}
    'domain.update': {}
}

export default mitt<Events>()
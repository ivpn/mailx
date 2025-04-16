import mitt from 'mitt'

// Define events
type Events = {
    'user.update': { email: string }
    'alias.create': {}
    'alias.update': {}
    'alias.delete': { id: string, catchAll: boolean }
    'totp.enable': {}
    'totp.disable': {}
    'recipient.create': {}
    'recipient.update': {}
    'recipient.verify': {}
    'recipient.delete': { id: string }
    'recipient.delete.error': { error: string }
    'recipient.reload': {}
}

export default mitt<Events>()
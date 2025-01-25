import mitt from 'mitt'

// Define events
type Events = {
    'user.update': { email: string }
    'alias.create': {}
    'alias.update': {}
    'alias.delete': { id: string }
    'totp.enable': {}
    'totp.disable': {}
    'recipient.create': {}
    'recipient.verify': {}
    'recipient.delete': { id: string }
}

export default mitt<Events>()
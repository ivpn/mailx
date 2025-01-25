import mitt from 'mitt'

// Define events
type Events = {
    'user.update': { email: string }
    'alias.create': {}
    'alias.update': {}
    'alias.delete': { id: string }
    'totp.enable': {}
    'totp.disable': {}
}

export default mitt<Events>()
import mitt from 'mitt'

// Define events
type Events = {
    'user.update': { email: string }
}

export default mitt<Events>()
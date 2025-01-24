import mitt from 'mitt'

// Define events
type Events = {
    'onUpdateEmail': { email: string }
}

export default mitt<Events>()
import mitt from 'mitt'

// Define events
type Events = {
    'alias.create': {}
}

export default mitt<Events>()
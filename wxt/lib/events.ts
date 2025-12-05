import mitt from 'mitt'
import { Alias } from '@/lib/types'

// Define events
type Events = {
    'alias.create': { alias: Alias }
}

export default mitt<Events>()
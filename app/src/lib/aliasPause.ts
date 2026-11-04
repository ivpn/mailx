import { onMounted, onUnmounted, ref } from 'vue'

export type PauseUnit = 'hours' | 'days' | 'weeks'

export const pauseValues = Array.from({ length: 10 }, (_, i) => i + 1)

export const pauseUnits: { value: PauseUnit, singular: string, plural: string }[] = [
    { value: 'hours', singular: 'Hour', plural: 'Hours' },
    { value: 'days', singular: 'Day', plural: 'Days' },
    { value: 'weeks', singular: 'Week', plural: 'Weeks' },
]

// Fixed lengths, the same as the server: a day is 24h and a week 7 × 24h.
const unitMs: Record<PauseUnit, number> = {
    hours: 60 * 60 * 1000,
    days: 24 * 60 * 60 * 1000,
    weeks: 7 * 24 * 60 * 60 * 1000,
}

// When a pause started at from would end. Display only; the server works out the real end time.
export const pauseEnd = (value: number, unit: PauseUnit, from: Date) => new Date(from.getTime() + value * unitMs[unit])

// One format for the row label, the disabled action's tooltip and the modal preview.
export const formatPausedUntil = (value: string | Date) => new Date(value).toLocaleString(undefined, { dateStyle: 'medium', timeStyle: 'short' })

// Matches ActiveStatus() on the server, like Rules.vue. An unknown status ('') isn't read-only.
const activeStatuses = ['active', 'grace_period']
export const isPauseReadOnly = (status: string) => status !== '' && !activeStatuses.includes(status)

// One clock for every row, ticking each minute while any component uses it, so a pause that
// ends while the page is open drops its label without a reload.
const now = ref(new Date())
let users = 0
let timer: ReturnType<typeof setInterval> | undefined

export const useNow = () => {
    onMounted(() => {
        if (users++ > 0) return
        now.value = new Date()
        timer = setInterval(() => {
            now.value = new Date()
        }, 60 * 1000)
    })

    onUnmounted(() => {
        if (--users > 0) return
        clearInterval(timer)
        timer = undefined
    })

    return now
}

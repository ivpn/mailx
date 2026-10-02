import { computed, ref } from 'vue'
import { announcementsApi } from '../api/announcements.ts'

// Announcements carry no id or timestamp, only a title with a trailing date
// ("... - 2026-09-01", see announcements.md). The newest one is listed first, so its
// title is the marker: the dot shows whenever the newest title differs from the last
// one the user saw on the Announcements page.
const STORAGE_KEY = 'announcements.seen'

// On a clean browser (nothing stored) only an announcement published in the last 48h
// counts as unread, so a fresh login doesn't light the dot for the whole backlog.
const CLEAN_BROWSER_UNREAD_WINDOW_MS = 48 * 60 * 60 * 1000

// The nav stays mounted for the whole session, so re-check when the tab regains focus,
// but no more often than this.
const REFRESH_INTERVAL_MS = 5 * 60 * 1000

const readSeen = (): string | null => {
    try {
        return localStorage.getItem(STORAGE_KEY)
    } catch {
        return null
    }
}

const latestTitle = ref<string | null>(null)
const seenTitle = ref<string | null>(readSeen())
let lastFetch = 0

const titleDate = (title: string): number | null => {
    const m = title.match(/(\d{4}-\d{2}-\d{2})\s*$/)
    if (!m) return null
    const ts = new Date(m[1]).getTime()
    return isNaN(ts) ? null : ts
}

const hasUnread = computed(() => {
    const latest = latestTitle.value
    if (!latest || latest === seenTitle.value) return false
    if (seenTitle.value !== null) return true

    const published = titleDate(latest)
    return published === null || Date.now() - published < CLEAN_BROWSER_UNREAD_WINDOW_MS
})

const setLatest = (list: { title: string }[]) => {
    latestTitle.value = Array.isArray(list) && list.length ? list[0].title : null
}

// A failed fetch leaves the previous state, it must never surface as a (false) dot.
const refresh = async () => {
    if (Date.now() - lastFetch < REFRESH_INTERVAL_MS) return
    lastFetch = Date.now()

    try {
        const response = await announcementsApi.getList()
        setLatest(response.data)
    } catch { }
}

// Called by the Announcements page with the list it just loaded.
const markSeen = (list: { title: string }[]) => {
    setLatest(list)
    lastFetch = Date.now()
    if (!latestTitle.value) return

    seenTitle.value = latestTitle.value
    try {
        localStorage.setItem(STORAGE_KEY, latestTitle.value)
    } catch { }
}

document.addEventListener('visibilitychange', () => {
    if (document.visibilityState === 'visible') refresh()
})

export const useAnnouncements = () => ({
    hasUnread,
    refresh,
    markSeen,
})

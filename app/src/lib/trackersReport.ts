// Decodes the removed trackers report that the API encodes in the fragment of the
// "View report" link in forwarded emails (base64url compact JSON). The fragment is never sent
// to the server, so the report is only ever read here in the browser.

export interface RemovedTracker {
    domain: string
    element: string
    reason: string
}

export interface TrackersReport {
    total: number
    trackers: RemovedTracker[]
}

export const elementLabels: Record<string, string> = {
    img: 'Image',
    srcset: 'Image srcset',
    background: 'Background image',
    link: 'Prefetch link',
}

export const reasonLabels: Record<string, string> = {
    domain: 'Known tracker domain',
    tiny: 'Tiny image',
    hidden: 'Hidden image',
    pattern: 'Tracking URL pattern',
}

const base64UrlToText = (encoded: string): string => {
    const base64 = encoded.replace(/-/g, '+').replace(/_/g, '/')
    const padded = base64 + '='.repeat((4 - (base64.length % 4)) % 4)
    const bytes = Uint8Array.from(atob(padded), (c) => c.charCodeAt(0))
    return new TextDecoder().decode(bytes)
}

// Returns null when the link is missing or malformed, e.g. truncated by a mail client.
export const decodeTrackersReport = (hash: string): TrackersReport | null => {
    const encoded = hash.replace(/^#/, '')
    if (!encoded) return null

    let data: unknown
    try {
        data = JSON.parse(base64UrlToText(encoded))
    } catch {
        return null
    }

    if (typeof data !== 'object' || data === null) return null
    const { n, t } = data as { n?: unknown, t?: unknown }
    if (typeof n !== 'number' || !Array.isArray(t)) return null

    const trackers = t
        .filter((item): item is Record<string, unknown> => typeof item === 'object' && item !== null)
        .map((item) => ({
            domain: typeof item.d === 'string' ? item.d : '',
            element: typeof item.e === 'string' ? item.e : '',
            reason: typeof item.r === 'string' ? item.r : '',
        }))

    return { total: Math.max(n, trackers.length), trackers }
}

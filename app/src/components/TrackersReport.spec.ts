import { afterEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App } from 'vue'
import TrackersReport from './TrackersReport.vue'
import { decodeTrackersReport } from '../lib/trackersReport.ts'

vi.mock('./Footer.vue', () => ({ default: { template: '<footer></footer>' } }))

// Produced by the API's trackers.EncodeReport for a Mailchimp open pixel
const mailchimpReport = 'eyJuIjoxLCJ0IjpbeyJkIjoic2hvcC51czEubGlzdC1tYW5hZ2UuY29tIiwiZSI6ImltZyIsInIiOiJkb21haW4ifV19'

const encode = (data: unknown) => {
    const bytes = new TextEncoder().encode(JSON.stringify(data))
    return btoa(String.fromCharCode(...bytes)).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '')
}

describe('decodeTrackersReport', () => {
    it('decodes the report the API encodes', () => {
        expect(decodeTrackersReport('#' + mailchimpReport)).toEqual({
            total: 1,
            trackers: [{ domain: 'shop.us1.list-manage.com', element: 'img', reason: 'domain' }],
        })
    })

    it('decodes non-ASCII domains and base64url characters', () => {
        const report = { n: 2, t: [{ d: 'tråck.example', e: 'srcset', r: 'pattern' }, { d: '???>>>', e: 'link', r: 'domain' }] }
        expect(decodeTrackersReport(encode(report))?.trackers.map((t) => t.domain)).toEqual(['tråck.example', '???>>>'])
    })

    it('keeps the total when the list was capped', () => {
        const report = { n: 120, t: [{ d: 'mailtrack.io', e: 'img', r: 'domain' }] }
        expect(decodeTrackersReport(encode(report))?.total).toBe(120)
    })

    it.each([
        '',
        '#',
        '#not-base64!',
        '#' + mailchimpReport.slice(0, 20),
        '#' + encode({ t: [] }),
        '#' + encode({ n: 1, t: 'x' }),
        '#' + encode(null),
    ])('returns null for an invalid report %s', (hash) => {
        expect(decodeTrackersReport(hash)).toBeNull()
    })

    it('ignores malformed entries', () => {
        const report = { n: 2, t: [null, 'x', { d: 5, e: 'img', r: 'tiny' }] }
        expect(decodeTrackersReport(encode(report))).toEqual({
            total: 2,
            trackers: [{ domain: '', element: 'img', reason: 'tiny' }],
        })
    })
})

describe('TrackersReport', () => {
    let app: App | undefined
    let root: HTMLElement

    const mount = async (hash: string) => {
        window.location.hash = hash
        root = document.createElement('div')
        document.body.appendChild(root)
        app = createApp(TrackersReport)
        app.mount(root)
        await nextTick()
    }

    afterEach(() => {
        app?.unmount()
        root.remove()
        window.location.hash = ''
    })

    it('renders the removed trackers', async () => {
        await mount(mailchimpReport)

        expect(root.textContent).toContain('Mailx removed 1 tracker from this email')
        const cells = [...root.querySelectorAll('tbody td')].map((td) => td.textContent?.trim())
        expect(cells).toEqual(['shop.us1.list-manage.com', 'Image', 'Known tracker domain'])
    })

    it('says when only the first trackers are listed', async () => {
        await mount(encode({ n: 60, t: [{ d: 'mailtrack.io', e: 'img', r: 'domain' }] }))

        expect(root.textContent).toContain('Mailx removed 60 trackers')
        expect(root.textContent).toContain('Showing the first 1 of 60 trackers.')
    })

    it('renders report values as text', async () => {
        await mount(encode({ n: 1, t: [{ d: '<img src=x onerror=alert(1)>', e: 'img', r: 'domain' }] }))

        expect(root.querySelector('tbody img')).toBeNull()
        expect(root.querySelector('tbody td')?.textContent).toBe('<img src=x onerror=alert(1)>')
    })

    it('shows an error for an invalid link', async () => {
        await mount('garbage')

        expect(root.querySelector('table')).toBeNull()
        expect(root.querySelector('.error')?.textContent).toContain('invalid or incomplete')
    })
})

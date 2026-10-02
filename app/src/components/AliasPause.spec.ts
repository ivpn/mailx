import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest'
import { createApp, nextTick, type App, type Component } from 'vue'
import { AxiosError, AxiosHeaders } from 'axios'
import Aliases from './Aliases.vue'
import Wildcards from './Wildcards.vue'
import { aliasApi } from '../api/alias.ts'
import { settingsApi } from '../api/settings.ts'
import { subscriptionApi } from '../api/subscription.ts'
import { formatPausedUntil, pauseEnd } from '../lib/aliasPause.ts'

// Preline's overlay is real here, as in RulesModals.spec.ts: it only opens a modal it has bound
// to a toggle, which a mocked overlay can't catch.
vi.mock('../api/alias.ts', () => ({
    aliasApi: { getList: vi.fn(), pause: vi.fn(), unpause: vi.fn() },
}))
vi.mock('../api/settings.ts', () => ({
    settingsApi: { getDefaults: vi.fn() },
}))
vi.mock('../api/subscription.ts', () => ({
    subscriptionApi: { get: vi.fn() },
}))
vi.mock('../lib/useBreakpoint.ts', async () => {
    const { ref } = await import('vue')
    return { useBreakpoint: () => ({ isDesktop: ref(true) }) }
})

let app: App | undefined
let root: HTMLElement

const wait = async (ms = 0) => {
    await new Promise((resolve) => setTimeout(resolve, ms))
    await nextTick()
}

const alias = {
    id: 'alias-1',
    created_at: '2026-09-01T10:00:00Z',
    deleted_at: null,
    name: 'promo@example.net',
    enabled: true,
    description: '',
    recipients: 'me@example.com',
    from_name: '',
    pinned: false,
    paused_until: null,
    origin: 0,
    is_custom_domain: false,
    is_domain_verified: null,
    is_domain_enabled: false,
    stats: { forwards: 0, blocks: 0, replies: 0, sends: 0 },
}

const isOpen = () => {
    const modal = document.getElementById('modal-alias-pause')!
    return !modal.classList.contains('hidden') && modal.classList.contains('open')
}

const valueSelect = () => document.querySelector<HTMLSelectElement>('#alias_pause_value')!
const unitSelect = () => document.querySelector<HTMLSelectElement>('#alias_pause_unit')!
const modalButton = (label: string) => [...document.querySelectorAll<HTMLButtonElement>('#modal-alias-pause footer button')].find((b) => b.textContent?.trim() === label)!
const pauseAction = () => root.querySelector<HTMLButtonElement>('[data-test="alias-pause"]')!

const choose = async (select: HTMLSelectElement, value: string) => {
    select.value = value
    select.dispatchEvent(new Event('change'))
    await nextTick()
}

const openModal = async () => {
    pauseAction().click()
    await wait(100)
}

const mount = async (status = 'active', page: Component = Aliases) => {
    vi.mocked(subscriptionApi.get).mockResolvedValue({ data: { status } } as any)

    root = document.createElement('div')
    document.body.appendChild(root)
    app = createApp(page)
    app.mount(root)
    await wait(50)
}

beforeEach(() => {
    vi.mocked(aliasApi.getList).mockReset().mockResolvedValue({ data: { aliases: [alias], total: 1 } } as any)
    vi.mocked(aliasApi.pause).mockReset()
    vi.mocked(settingsApi.getDefaults).mockResolvedValue({ data: { id: 'settings-1', recipients: [], domains: [], custom_domains: [] } } as any)
})

afterEach(() => {
    app?.unmount()
    root.remove()
    document.querySelectorAll('[data-hs-overlay-backdrop-template], .hs-overlay-backdrop').forEach((el) => el.remove())
    document.body.style.overflow = ''
})

describe('AliasPause', () => {
    it('opens from the row with 1 Day selected and a preview', async () => {
        await mount()
        expect(isOpen()).toBe(false)

        await openModal()

        expect(isOpen()).toBe(true)
        expect(document.querySelector('#modal-alias-pause h4')?.textContent).toBe('PAUSE ALIAS')
        expect(document.querySelector('#modal-alias-pause article')?.textContent).toContain('Emails forwarded to promo@example.net will be discarded until the pause ends.')
        expect(valueSelect().value).toBe('1')
        expect(unitSelect().value).toBe('days')
        expect(unitSelect().selectedOptions[0].textContent).toBe('Day')
        expect(document.querySelector('[data-test="pause-preview"]')?.textContent)
            .toBe('Paused until ' + formatPausedUntil(pauseEnd(1, 'days', new Date())))
    })

    it('sends the chosen duration, closes and reloads the list', async () => {
        vi.mocked(aliasApi.pause).mockResolvedValue({ data: { message: 'ok', paused_until: '2026-10-02T13:00:00Z' } } as any)
        await mount()
        await openModal()

        await choose(valueSelect(), '3')
        await choose(unitSelect(), 'hours')
        expect(unitSelect().selectedOptions[0].textContent).toBe('Hours')
        expect(document.querySelector('[data-test="pause-preview"]')?.textContent)
            .toBe('Paused until ' + formatPausedUntil(pauseEnd(3, 'hours', new Date())))

        const listCalls = vi.mocked(aliasApi.getList).mock.calls.length
        modalButton('Pause').click()
        await wait(100)

        expect(aliasApi.pause).toHaveBeenCalledWith('alias-1', 3, 'hours')
        expect(isOpen()).toBe(false)
        expect(vi.mocked(aliasApi.getList).mock.calls.length).toBeGreaterThan(listCalls)
    })

    it('resets the selects when reopened', async () => {
        await mount()
        await openModal()
        await choose(valueSelect(), '7')
        await choose(unitSelect(), 'weeks')

        modalButton('Cancel').click()
        await wait(100)
        expect(isOpen()).toBe(false)

        await openModal()
        expect(valueSelect().value).toBe('1')
        expect(unitSelect().value).toBe('days')
    })

    it('keeps the modal open and shows a server error', async () => {
        const config = { headers: new AxiosHeaders() }
        const message = 'Your subscription is not active. Please renew to pause aliases.'
        vi.mocked(aliasApi.pause).mockRejectedValue(new AxiosError('Request failed', 'ERR_BAD_REQUEST', config, undefined, {
            status: 400, statusText: 'Bad Request', headers: {}, config, data: { error: message },
        }))
        await mount()
        await openModal()

        modalButton('Pause').click()
        await wait(100)

        expect(isOpen()).toBe(true)
        expect(document.querySelector('#modal-alias-pause footer .error')?.textContent).toBe('Error: ' + message)
        expect(modalButton('Pause').disabled).toBe(false)
    })

    it('disables the button while the request is running', async () => {
        let resolve: (value: unknown) => void = () => {}
        vi.mocked(aliasApi.pause).mockReturnValue(new Promise((r) => { resolve = r }) as any)
        await mount()
        await openModal()

        modalButton('Pause').click()
        await nextTick()
        expect(modalButton('Pause').disabled).toBe(true)

        resolve({ data: {} })
        await wait(100)
    })

    it('works from the Wildcard list', async () => {
        const wildcard = { ...alias, name: '*+news@example.net', wildcard: true }
        vi.mocked(aliasApi.getList).mockResolvedValue({ data: { aliases: [wildcard], total: 1 } } as any)
        vi.mocked(aliasApi.pause).mockResolvedValue({ data: {} } as any)
        await mount('active', Wildcards)
        await openModal()

        expect(isOpen()).toBe(true)
        expect(document.querySelector('#modal-alias-pause article')?.textContent).toContain('*+news@example.net')

        modalButton('Pause').click()
        await wait(100)
        expect(aliasApi.pause).toHaveBeenCalledWith('alias-1', 1, 'days')
    })

    it('is disabled with a tooltip while the subscription is limited', async () => {
        await mount('limited_access')

        expect(pauseAction().disabled).toBe(true)
        expect(pauseAction().parentElement?.querySelector('[role="tooltip"]')?.textContent?.trim()).toBe('Renew your subscription to pause aliases.')
    })

    it.each(['active', 'grace_period'])('is enabled while the subscription is %s', async (status) => {
        await mount(status)

        expect(pauseAction().disabled).toBe(false)
        expect(pauseAction().parentElement?.querySelector('[role="tooltip"]')).toBeNull()
    })
})
